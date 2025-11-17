package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/nyxstack/color"
)

// Command represents a command in the CLI application
type Command struct {
	name        string
	description string
	args        []Argument
	flags       *FlagSet // All flags (automatically inherit to children)
	subcommands map[string]*Command
	parent      *Command
	action      interface{}
	hidden      bool

	// Lifecycle hooks
	persistentPreRun  func(context.Context, *Command) error
	preRun            func(context.Context, *Command) error
	postRun           func(context.Context, *Command) error
	persistentPostRun func(context.Context, *Command) error

	// Help system
	helpEnabled bool
	helpFlag    string
	helpShort   string
}

// Getter methods (public API)
func (c *Command) GetName() string {
	return c.name
}

func (c *Command) GetDescription() string {
	return c.description
}

func (c *Command) GetParent() *Command {
	return c.parent
}

func (c *Command) GetCommands() map[string]*Command {
	return c.subcommands
}

func (c *Command) GetArgs() []Argument {
	return c.args
}

// Cmd creates a new command with the given name
func Cmd(name string) *Command {
	return &Command{
		name:        name,
		flags:       NewFlagSet(),
		subcommands: make(map[string]*Command),
		helpEnabled: true,
		helpFlag:    "help",
		helpShort:   "h",
	}
}

// Root creates a new root command (convenience function)
func Root(name string) *Command {
	return Cmd(name)
}

// Hidden marks the command as hidden from help output
func (c *Command) Hidden() *Command {
	c.hidden = true
	return c
}

// Show marks the command as visible in help output (default)
func (c *Command) Show() *Command {
	c.hidden = false
	return c
}

// IsHidden returns whether the command is hidden
func (c *Command) IsHidden() bool {
	return c.hidden
}

// Description sets the command description
func (c *Command) Description(desc string) *Command {
	c.description = desc
	return c
}

// Flag adds a typed flag to the command using reflection
func (c *Command) Flag(ptr interface{}, name, shorthand string, defaultValue interface{}, usage string) *Command {
	c.flags.Add(ptr, name, shorthand, defaultValue, usage)
	return c
}

// FlagRequired adds a required flag to the command
func (c *Command) FlagRequired(ptr interface{}, name, shorthand string, defaultValue interface{}, usage string) *Command {
	c.flags.Add(ptr, name, shorthand, defaultValue, usage)
	// Mark the flag as required
	if flag := c.flags.GetFlag(name); flag != nil {
		flag.required = true
	}
	return c
}

// FlagHidden adds a hidden flag to the command
func (c *Command) FlagHidden(ptr interface{}, name, shorthand string, defaultValue interface{}, usage string) *Command {
	c.flags.Add(ptr, name, shorthand, defaultValue, usage)
	// Mark the flag as hidden
	if flag := c.flags.GetFlag(name); flag != nil {
		flag.hidden = true
	}
	return c
}

// Flags binds struct fields as flags using struct tags
func (c *Command) Flags(structPtr interface{}) *Command {
	c.flags.BindStruct(structPtr)
	return c
}

// Arg adds a positional argument to the command
func (c *Command) Arg(name, description string, required bool) *Command {
	c.args = append(c.args, Argument{
		Name:        name,
		Description: description,
		Required:    required,
	})
	return c
}

// Action sets the function to execute when this command is run
func (c *Command) Action(fn interface{}) *Command {
	c.action = fn
	return c
}

// PersistentPreRun sets a function to run before this command and all subcommands (inherits to children)
func (c *Command) PersistentPreRun(fn func(context.Context, *Command) error) *Command {
	c.persistentPreRun = fn
	return c
}

// PreRun sets a function to run before this command's action (command-specific)
func (c *Command) PreRun(fn func(context.Context, *Command) error) *Command {
	c.preRun = fn
	return c
}

// PostRun sets a function to run after this command's action (command-specific)
func (c *Command) PostRun(fn func(context.Context, *Command) error) *Command {
	c.postRun = fn
	return c
}

// PersistentPostRun sets a function to run after this command and all subcommands (inherits to children)
func (c *Command) PersistentPostRun(fn func(context.Context, *Command) error) *Command {
	c.persistentPostRun = fn
	return c
}

// AddCommand adds a subcommand
func (c *Command) AddCommand(cmd *Command) *Command {
	cmd.parent = c
	c.subcommands[cmd.name] = cmd
	return c
}

// getAllFlags returns all flags including inherited from ancestors
func (c *Command) getAllFlags() []*Flag {
	var allFlags []*Flag
	seen := make(map[string]bool) // Track by primary name to avoid duplicates

	// Collect flags from ancestors (root to current)
	var ancestors []*Command
	current := c
	for current != nil {
		ancestors = append([]*Command{current}, ancestors...)
		current = current.parent
	}

	// Add flags from current to root (child flags shadow parent flags)
	for i := len(ancestors) - 1; i >= 0; i-- {
		cmd := ancestors[i]
		for _, flag := range cmd.flags.GetFlags() {
			primaryName := flag.PrimaryName()
			if !seen[primaryName] {
				allFlags = append(allFlags, flag)
				seen[primaryName] = true
			}
		}
	}

	return allFlags
}

// getCommandPath returns the full command path from root to this command
func (c *Command) getCommandPath() string {
	if c.parent == nil {
		return c.name
	}
	return c.parent.getCommandPath() + " " + c.name
}

// DisableHelp disables the automatic help functionality
func (c *Command) DisableHelp() *Command {
	c.helpEnabled = false
	return c
}

// EnableHelp enables the automatic help functionality (default)
func (c *Command) EnableHelp() *Command {
	c.helpEnabled = true
	return c
}

// SetHelpFlag sets custom help flag names (default: "help" and "h")
func (c *Command) SetHelpFlag(long, short string) *Command {
	c.helpFlag = long
	c.helpShort = short
	return c
}

// IsHelpEnabled returns whether help is enabled for this command
func (c *Command) IsHelpEnabled() bool {
	return c.helpEnabled
}

// showHelp displays help information for the command
func (c *Command) showHelp() {
	// Build the full command path for usage
	commandPath := c.getCommandPath()
	fmt.Printf("%s: %s", color.Bold+"Usage"+color.Reset, commandPath)

	// Show subcommands indicator first
	if len(c.subcommands) > 0 {
		fmt.Printf(" %s", color.Cyan+"[command]"+color.Reset)
	}

	// Show arguments after subcommands
	if len(c.args) > 0 {
		argList := []string{}
		for _, arg := range c.args {
			if arg.Required {
				argList = append(argList, fmt.Sprintf("<%s>", arg.Name))
			} else {
				argList = append(argList, fmt.Sprintf("[%s]", arg.Name))
			}
		}
		if len(argList) > 0 {
			fmt.Printf(" %s", color.Yellow+strings.Join(argList, " ")+color.Reset)
		}
	}

	// Show flags indicator if any flags exist (local or inherited)
	allFlags := c.getAllFlags()
	if len(allFlags) > 0 || c.helpEnabled {
		fmt.Printf(" %s", color.Dim+"[flags...]"+color.Reset)
	}

	fmt.Println()

	if c.description != "" {
		fmt.Printf("\n%s\n", c.description)
	}

	// Show arguments with descriptions
	if len(c.args) > 0 {
		fmt.Printf("\n%s:\n", color.Bold+"Arguments"+color.Reset)
		for _, arg := range c.args {
			required := ""
			if arg.Required {
				required = " " + color.Red + "(required)" + color.Reset
			} else {
				required = " " + color.Dim + "(optional)" + color.Reset
			}
			fmt.Printf("  %-15s %s%s\n", color.Yellow+arg.Name+color.Reset, arg.Description, required)
		}
	}

	// Show all flags (local and inherited)
	if len(allFlags) > 0 {
		fmt.Printf("\n%s:\n", color.Bold+"Flags"+color.Reset)

		// Track displayed flags by primary name to avoid duplicates
		displayed := make(map[string]bool)

		// Display flags (child command's local flags take precedence over inherited)
		for _, flag := range allFlags {
			primaryName := flag.PrimaryName()

			if !displayed[primaryName] {
				// Determine if this is an inherited flag
				isLocal := false
				for _, localFlag := range c.flags.GetFlags() {
					if localFlag.PrimaryName() == primaryName {
						isLocal = true
						break
					}
				}

				suffix := ""
				if !isLocal && c.parent != nil {
					suffix = color.Dim + " (inherited)" + color.Reset
				}

				c.displayFlag(flag, suffix)
				displayed[primaryName] = true
			}
		}

		// Add help flag if enabled
		if c.helpEnabled {
			helpNames := fmt.Sprintf("%s, %s", color.Green+fmt.Sprintf("-%s", c.helpShort)+color.Reset, color.Green+fmt.Sprintf("--%s", c.helpFlag)+color.Reset)
			fmt.Printf("  %-30s %s\n", helpNames, "Show help information")
		}
	} else if c.helpEnabled {
		// Show help flag even if no other flags
		fmt.Printf("\n%s:\n", color.Bold+"Flags"+color.Reset)
		helpNames := fmt.Sprintf("%s, %s", color.Green+fmt.Sprintf("-%s", c.helpShort)+color.Reset, color.Green+fmt.Sprintf("--%s", c.helpFlag)+color.Reset)
		fmt.Printf("  %-30s %s\n", helpNames, "Show help information")
	}

	// Show subcommands
	if len(c.subcommands) > 0 {
		// Count visible subcommands
		visibleCount := 0
		for _, cmd := range c.subcommands {
			if !cmd.IsHidden() {
				visibleCount++
			}
		}

		if visibleCount > 0 {
			fmt.Printf("\n%s:\n", color.Bold+"Commands"+color.Reset)
			for name, cmd := range c.subcommands {
				if !cmd.IsHidden() {
					fmt.Printf("  %-15s %s\n", color.Cyan+name+color.Reset, cmd.description)
				}
			}

			// Show help command if enabled
			if c.helpEnabled {
				fmt.Printf("\n%s \"%s [command] %s\" %s\n",
					color.Dim+"Use"+color.Reset,
					c.name,
					color.Green+"--"+c.helpFlag+color.Reset,
					color.Dim+"for more information about a command."+color.Reset)
			}
		}
	}
}

// displayFlag formats and displays a single flag
func (c *Command) displayFlag(flag *Flag, suffix string) {
	names := color.Green + fmt.Sprintf("--%s", flag.PrimaryName()) + color.Reset
	if flag.ShortName() != "" {
		names = fmt.Sprintf("%s, %s", color.Green+fmt.Sprintf("-%s", flag.ShortName())+color.Reset, names)
	}

	defaultInfo := ""
	if flag.GetDefault() != nil {
		defaultInfo = color.Dim + fmt.Sprintf(" (default: %v)", flag.GetDefault()) + color.Reset
	}

	fmt.Printf("  %-30s %s%s%s\n", names, flag.GetUsage(), defaultInfo, suffix)
}

// ShowHelp displays help information (public API)
func (c *Command) ShowHelp() {
	c.showHelp()
}
