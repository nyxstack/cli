package cli

import (
	"context"
	"fmt"
)

// BashCompletion implements bash shell completion
type BashCompletion struct{}

func (b *BashCompletion) GetCompletions(cmd *Command, args []string) []string {
	return getCompletionWords(cmd)
}

func (b *BashCompletion) Register(cmd *Command) {
	bashCmd := Cmd("__bashcomplete").
		Description("Bash completion helper").
		Hidden().
		Action(func(ctx context.Context, bashCommand *Command) error {
			targetCmd := bashCommand.GetParent()
			// For completion, we don't need args since we complete the parent
			words := b.GetCompletions(targetCmd, nil)

			for _, word := range words {
				fmt.Println(word)
			}
			return nil
		})

	cmd.AddCommand(bashCmd)

	// Recursively register for all subcommands
	for _, subcmd := range cmd.GetCommands() {
		if !subcmd.IsHidden() {
			b.Register(subcmd)
		}
	}
}

func (b *BashCompletion) GenerateScript(cmd *Command) string {
	cmdName := cmd.GetName()

	script := fmt.Sprintf(`# Bash completion script for %s
# Source this file to enable bash completion:
#   source <(%s __bashcomplete)

_%s_completion() {
    local cur prev words cword
    _init_completion || return

    # Get the full command path
    local cmd_path="${COMP_WORDS[0]}"
    for ((i=1; i < COMP_CWORD; i++)); do
        local word="${COMP_WORDS[i]}"
        if [[ "$word" != -* ]]; then
            cmd_path="$cmd_path $word"
        fi
    done

    # Get completions from the command
    local completions=$($cmd_path __bashcomplete 2>/dev/null)
    
    # Generate reply
    COMPREPLY=($(compgen -W "$completions" -- "$cur"))
}

complete -F _%s_completion %s
`, cmdName, cmdName, cmdName, cmdName, cmdName)

	return script
}
