package cli

import (
	"context"
	"fmt"
)

// PowerShellCompletion implements PowerShell completion
type PowerShellCompletion struct{}

func (p *PowerShellCompletion) GetCompletions(cmd *Command, args []string) []string {
	return getCompletionWords(cmd)
}

func (p *PowerShellCompletion) Register(cmd *Command) {
	psCmd := Cmd("__powershellcomplete").
		Description("PowerShell completion helper").
		Hidden().
		Action(func(ctx context.Context, psCommand *Command) error {
			targetCmd := psCommand.GetParent()
			// For completion, we don't need args since we complete the parent
			words := p.GetCompletions(targetCmd, nil)

			for _, word := range words {
				fmt.Println(word)
			}
			return nil
		})

	cmd.AddCommand(psCmd)

	// Recursively register for all subcommands
	for _, subcmd := range cmd.GetCommands() {
		if !subcmd.IsHidden() {
			p.Register(subcmd)
		}
	}
}

func (p *PowerShellCompletion) GenerateScript(cmd *Command) string {
	cmdName := cmd.GetName()

	script := fmt.Sprintf(`# PowerShell completion script for %s
# Add this to your PowerShell profile:
#   %s __powershellcomplete | Out-String | Invoke-Expression

Register-ArgumentCompleter -Native -CommandName %s -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    
    $cmdPath = $commandAst.CommandElements[0].Value
    
    # Build command path from AST
    for ($i = 1; $i -lt $commandAst.CommandElements.Count; $i++) {
        $element = $commandAst.CommandElements[$i].Value
        if (-not $element.StartsWith('-')) {
            $cmdPath = "$cmdPath $element"
        }
    }
    
    # Get completions
    $completions = & $cmdPath __powershellcomplete 2>$null
    
    $completions | ForEach-Object {
        [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
    }
}
`, cmdName, cmdName, cmdName)

	return script
}
