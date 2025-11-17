package cli

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// Execute runs the command with os.Args
func (c *Command) Execute() error {
	return c.ExecuteContext(context.Background())
}

// ExecuteContext runs the command with a context
func (c *Command) ExecuteContext(ctx context.Context) error {
	// Use os.Args[1:] (skip program name)
	args := os.Args[1:]
	return c.execute(ctx, args)
}

// ExecuteWithArgs runs the command with custom arguments (useful for testing)
func (c *Command) ExecuteWithArgs(args []string) error {
	return c.execute(context.Background(), args)
}

// execute is the internal execution logic
func (c *Command) execute(ctx context.Context, args []string) error {
	// Check for help flag first - find which command needs help
	if c.helpEnabled {
		for _, arg := range args {
			if arg == "--"+c.helpFlag || arg == "-"+c.helpShort {
				// Find which command the help is for
				targetCmd := c
				for _, a := range args {
					if !strings.HasPrefix(a, "-") {
						if cmd, exists := c.subcommands[a]; exists {
							targetCmd = cmd
							break
						}
					}
				}
				targetCmd.showHelp()
				return nil
			}
		}
	}

	// First, find if there's a subcommand in the args
	subcommandIndex := -1
	var subcmd *Command

	for i, arg := range args {
		// Skip flags (all flags start with - or --)
		if strings.HasPrefix(arg, "-") {
			continue
		}

		// Check if this is a known subcommand
		if cmd, exists := c.subcommands[arg]; exists {
			subcommandIndex = i
			subcmd = cmd
			break
		}

		// If it's not a flag and not a subcommand:
		// - If we have subcommands defined, this is an unknown command error
		// - Otherwise, it's an argument - stop looking for subcommands
		if len(c.subcommands) > 0 {
			// We have subcommands but this doesn't match any
			return &CommandNotFoundError{
				Name: arg,
				Cmd:  c,
			}
		}
		break
	}

	// If we found a subcommand, delegate to it
	if subcommandIndex >= 0 {
		// Parse flags that appear before the subcommand (belong to current command)
		beforeSubcmd := args[:subcommandIndex]
		afterSubcmd := args[subcommandIndex+1:]

		// Parse current command's flags from args before subcommand
		if len(beforeSubcmd) > 0 {
			allFlags := c.getAllFlags()
			tempFS := NewFlagSet()
			for _, flag := range allFlags {
				tempFS.flags = append(tempFS.flags, flag)
			}

			_, err := tempFS.Parse(beforeSubcmd)
			if err != nil {
				return &FlagError{
					Flag: "",
					Msg:  err.Error(),
					Cmd:  c,
				}
			}
		}

		// Execute subcommand with args after the subcommand name
		return subcmd.execute(ctx, afterSubcmd)
	}

	// No subcommand found, parse all flags and execute this command
	allFlags := c.getAllFlags()
	tempFS := NewFlagSet()
	for _, flag := range allFlags {
		tempFS.flags = append(tempFS.flags, flag)
	}

	remaining, err := tempFS.Parse(args)
	if err != nil {
		return &FlagError{
			Flag: "",
			Msg:  err.Error(),
			Cmd:  c,
		}
	}

	// Validate required flags
	for _, flag := range allFlags {
		if flag.IsRequired() && !flag.IsSet() {
			return &FlagError{
				Flag: flag.names[0],
				Msg:  "required flag not set",
				Cmd:  c,
			}
		}
	}

	// Validate argument count
	nonFlagArgs := remaining
	expectedArgs := len(c.args)

	// Check if action is variadic
	isVariadic := false
	if c.action != nil {
		actionValue := reflect.ValueOf(c.action)
		actionType := actionValue.Type()
		isVariadic = actionType.IsVariadic()
	}

	// Check if we have too many arguments (skip check for variadic)
	if !isVariadic && len(nonFlagArgs) > expectedArgs {
		return &ArgumentError{
			Arg: "",
			Msg: fmt.Sprintf("too many arguments: expected %d, got %d", expectedArgs, len(nonFlagArgs)),
			Cmd: c,
		}
	}

	// Execute this command's action
	return c.executeAction(ctx, remaining)
} // executeAction executes the command's action with lifecycle hooks
func (c *Command) executeAction(ctx context.Context, args []string) error {
	// Run PersistentPreRun hooks (from root to current)
	var ancestors []*Command
	current := c
	for current != nil {
		ancestors = append([]*Command{current}, ancestors...)
		current = current.parent
	}

	for _, cmd := range ancestors {
		if cmd.persistentPreRun != nil {
			if err := cmd.persistentPreRun(ctx, c); err != nil {
				// Still run post hooks on error
				c.runPostHooks(ctx)
				return err
			}
		}
	}

	// Run PreRun hook
	if c.preRun != nil {
		if err := c.preRun(ctx, c); err != nil {
			// Run post hooks even on PreRun error
			c.runPostHooks(ctx)
			return err
		}
	}

	// Execute action
	var actionErr error
	if c.action != nil {
		actionErr = c.callAction(ctx, args)
	}

	// Always run post hooks (even if action failed)
	c.runPostHooks(ctx)

	return actionErr
}

// runPostHooks executes PostRun and PersistentPostRun hooks
func (c *Command) runPostHooks(ctx context.Context) {
	// Run PostRun hook
	if c.postRun != nil {
		c.postRun(ctx, c) // Ignore errors in PostRun for now
	}

	// Run PersistentPostRun hooks (from current to root)
	current := c
	for current != nil {
		if current.persistentPostRun != nil {
			current.persistentPostRun(ctx, c) // Ignore errors
		}
		current = current.parent
	}
}

// callAction invokes the action function with proper arguments
func (c *Command) callAction(ctx context.Context, args []string) error {
	actionValue := reflect.ValueOf(c.action)
	actionType := actionValue.Type()

	// Build argument list
	callArgs := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(c),
	}

	// Check if function is variadic
	isVariadic := actionType.IsVariadic()
	numParams := actionType.NumIn() - 2 // Subtract ctx and cmd

	if isVariadic {
		// For variadic functions, handle specially
		// Add non-variadic arguments first (all params except the last variadic one)
		numFixed := numParams - 1

		for i := 0; i < numFixed; i++ {
			if i >= len(args) {
				// Check if this argument is required
				if i < len(c.args) && c.args[i].Required {
					return &ArgumentError{
						Arg: c.args[i].Name,
						Msg: "required argument missing",
						Cmd: c,
					}
				}
				// Use zero value for optional arguments
				callArgs = append(callArgs, reflect.Zero(actionType.In(i+2)))
				continue
			}

			// Convert string argument to expected type
			argType := actionType.In(i + 2)
			argValue, err := convertArgument(args[i], argType)
			if err != nil {
				argName := ""
				if i < len(c.args) {
					argName = c.args[i].Name
				}
				return &ArgumentError{
					Arg: argName,
					Msg: err.Error(),
					Cmd: c,
				}
			}

			callArgs = append(callArgs, argValue)
		}

		// Get the element type of the variadic parameter
		sliceType := actionType.In(actionType.NumIn() - 1)
		elemType := sliceType.Elem()

		// Build a slice for variadic parameters
		variadicCount := len(args) - numFixed
		variadicSlice := reflect.MakeSlice(sliceType, variadicCount, variadicCount)

		for i := 0; i < variadicCount; i++ {
			argValue, err := convertArgument(args[numFixed+i], elemType)
			if err != nil {
				return &ArgumentError{
					Arg: "",
					Msg: err.Error(),
					Cmd: c,
				}
			}
			variadicSlice.Index(i).Set(argValue)
		}

		// Append the slice as a single argument
		callArgs = append(callArgs, variadicSlice)
	} else {
		// Non-variadic function
		for i := 0; i < numParams; i++ {
			if i >= len(args) {
				// Check if this argument is required
				if i < len(c.args) && c.args[i].Required {
					return &ArgumentError{
						Arg: c.args[i].Name,
						Msg: "required argument missing",
						Cmd: c,
					}
				}
				// Use zero value for optional arguments
				callArgs = append(callArgs, reflect.Zero(actionType.In(i+2)))
				continue
			}

			// Convert string argument to expected type
			argType := actionType.In(i + 2)
			argValue, err := convertArgument(args[i], argType)
			if err != nil {
				argName := ""
				if i < len(c.args) {
					argName = c.args[i].Name
				}
				return &ArgumentError{
					Arg: argName,
					Msg: err.Error(),
					Cmd: c,
				}
			}

			callArgs = append(callArgs, argValue)
		}
	}

	// Call the action function
	var results []reflect.Value
	if isVariadic {
		// For variadic functions, use CallSlice with the slice as the last argument
		// CallSlice will expand the slice elements into individual variadic parameters
		results = actionValue.CallSlice(callArgs)
	} else {
		results = actionValue.Call(callArgs)
	}

	// Check if there's an error return value
	if len(results) > 0 && !results[0].IsNil() {
		return results[0].Interface().(error)
	}

	return nil
}

// convertArgument converts a string argument to the target type
func convertArgument(arg string, targetType reflect.Type) (reflect.Value, error) {
	// Use the same conversion logic as flag parsing
	tempFS := NewFlagSet()

	// Create a temporary variable of the target type
	tempVar := reflect.New(targetType).Elem()

	// Create a temporary flag
	tempFlag := Flag{
		names:    []string{"temp"},
		flagType: targetType,
		value:    tempVar,
	}

	// Parse the value
	if err := tempFS.setValue(&tempFlag, arg); err != nil {
		return reflect.Value{}, err
	}

	return tempVar, nil
}
