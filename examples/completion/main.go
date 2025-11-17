// Shell Completion Example
//
// WHAT: Demonstrates automatic shell completion for bash, zsh, fish, and PowerShell.
//
// WHY: Shell completion improves UX by:
//   - Suggesting available commands and subcommands
//   - Showing available flags and their shortcuts
//   - Reducing typos and documentation lookups
//   - Supporting multiple shell environments
//
// HOW IT WORKS:
//   - cli.AddCompletion(rootCmd) registers __<shell>complete commands
//   - Each command gets its own completion command
//   - Hidden commands (like debug) are excluded from completions
//   - Flags from parent commands are inherited
//
// USAGE:
//
//	go run main.go deploy <TAB>               # Complete environment argument
//	go run main.go database <TAB>             # Show 'migrate' subcommand
//	go run main.go --<TAB>                    # Show available flags
//	go run main.go database migrate --<TAB>   # Show flags (including inherited timeout)
//
// SETUP (choose your shell):
//
//	# Bash
//	source <(./myapp __bashcomplete)
//
//	# Zsh
//	source <(./myapp __zshcomplete)
//
//	# Fish
//	./myapp __fishcomplete | source
//
//	# PowerShell
//	./myapp __powershellcomplete | Out-String | Invoke-Expression
//
// KEY CONCEPTS:
//   - AddCompletion: Registers all shell completion commands
//   - Context-aware: Each command knows its subcommands and flags
//   - Hidden commands: Excluded from completion suggestions
//   - Flag inheritance: Child completions include parent flags
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nyxstack/cli"
)

var (
	verbose bool
	timeout int
)

var rootCmd = cli.Root("myapp").
	Description("Application with shell completion support").
	Flag(&verbose, "verbose", "v", false, "Enable verbose output")

var deployCmd = cli.Cmd("deploy").
	Description("Deploy the application").
	Arg("environment", "Target environment", true).
	Action(func(ctx context.Context, cmd *cli.Command, env string) error {
		fmt.Printf("Deploying to %s\n", env)
		return nil
	})

var dbCmd = cli.Cmd("database").
	Description("Database operations").
	Flag(&timeout, "timeout", "t", 30, "Operation timeout in seconds")

var migrateCmd = cli.Cmd("migrate").
	Description("Run database migrations").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Printf("Running migrations (timeout: %ds)\n", timeout)
		return nil
	})

var debugCmd = cli.Cmd("debug").
	Description("Debug command").
	Hidden().
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("Debug info...")
		return nil
	})

func init() {
	dbCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(debugCmd)

	// Enable shell completion
	cli.AddCompletion(rootCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
