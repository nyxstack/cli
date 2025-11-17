package cli

import (
	"context"
	"fmt"
)

// FishCompletion implements fish shell completion
type FishCompletion struct{}

func (f *FishCompletion) GetCompletions(cmd *Command, args []string) []string {
	return getCompletionWords(cmd)
}

func (f *FishCompletion) Register(cmd *Command) {
	fishCmd := Cmd("__fishcomplete").
		Description("Fish completion helper").
		Hidden().
		Action(func(ctx context.Context, fishCommand *Command) error {
			targetCmd := fishCommand.GetParent()
			// For completion, we don't need args since we complete the parent
			words := f.GetCompletions(targetCmd, nil)

			for _, word := range words {
				fmt.Println(word)
			}
			return nil
		})

	cmd.AddCommand(fishCmd)

	// Recursively register for all subcommands
	for _, subcmd := range cmd.GetCommands() {
		if !subcmd.IsHidden() {
			f.Register(subcmd)
		}
	}
}

func (f *FishCompletion) GenerateScript(cmd *Command) string {
	cmdName := cmd.GetName()

	script := fmt.Sprintf(`# Fish completion script for %s
# Save this to ~/.config/fish/completions/%s.fish

function __%s_complete
    set -l cmd_path (commandline -opc)
    $cmd_path __fishcomplete 2>/dev/null
end

complete -c %s -f -a "(__%s_complete)"
`, cmdName, cmdName, cmdName, cmdName, cmdName)

	return script
}
