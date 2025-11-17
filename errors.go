package cli

import "fmt"

// CommandNotFoundError indicates a subcommand was not found
type CommandNotFoundError struct {
	Name string
	Cmd  *Command
}

func (e *CommandNotFoundError) Error() string {
	return fmt.Sprintf("unknown command '%s' for '%s'", e.Name, e.Cmd.getCommandPath())
}

// ArgumentError indicates an argument validation error
type ArgumentError struct {
	Arg string
	Msg string
	Cmd *Command
}

func (e *ArgumentError) Error() string {
	return fmt.Sprintf("argument '%s': %s", e.Arg, e.Msg)
}

// FlagError indicates a flag parsing or validation error
type FlagError struct {
	Flag string
	Msg  string
	Cmd  *Command
}

func (e *FlagError) Error() string {
	return fmt.Sprintf("flag '%s': %s", e.Flag, e.Msg)
}
