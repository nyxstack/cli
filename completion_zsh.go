package cli

import (
	"context"
	"fmt"
)

// ZshCompletion implements zsh shell completion
type ZshCompletion struct{}

func (z *ZshCompletion) GetCompletions(cmd *Command, args []string) []string {
	return getCompletionWords(cmd)
}

func (z *ZshCompletion) Register(cmd *Command) {
	zshCmd := Cmd("__zshcomplete").
		Description("Zsh completion helper").
		Hidden().
		Action(func(ctx context.Context, zshCommand *Command) error {
			targetCmd := zshCommand.GetParent()
			// For completion, we don't need args since we complete the parent
			words := z.GetCompletions(targetCmd, nil)

			for _, word := range words {
				fmt.Println(word)
			}
			return nil
		})

	cmd.AddCommand(zshCmd)

	// Recursively register for all subcommands
	for _, subcmd := range cmd.GetCommands() {
		if !subcmd.IsHidden() {
			z.Register(subcmd)
		}
	}
}

func (z *ZshCompletion) GenerateScript(cmd *Command) string {
	cmdName := cmd.GetName()

	script := fmt.Sprintf(`#compdef %s

# Zsh completion script for %s
# Add this to your fpath and reload completions:
#   %s __zshcomplete > ~/.zsh/completions/_%s
#   compinit

_%s() {
    local -a completions
    local cmd_path="${words[1]}"
    
    # Build command path from words
    for ((i=2; i < CURRENT; i++)); do
        if [[ "${words[i]}" != -* ]]; then
            cmd_path="$cmd_path ${words[i]}"
        fi
    done
    
    # Get completions
    completions=(${(f)"$($cmd_path __zshcomplete 2>/dev/null)"})
    
    _describe '%s' completions
}

_%s "$@"
`, cmdName, cmdName, cmdName, cmdName, cmdName, cmdName, cmdName)

	return script
}
