// Lifecycle Hooks Example
//
// WHAT: Demonstrates the execution order of lifecycle hooks in the CLI framework.
//
// WHY: Understanding hook execution order is crucial for:
//   - Proper initialization and cleanup
//   - Validation before action execution
//   - Shared setup/teardown logic across multiple commands
//   - Resource management (database connections, file handles, etc.)
//
// EXECUTION ORDER:
//  1. PersistentPreRun  (parent, then child) - Global initialization
//  2. PreRun            (child only)         - Command-specific setup/validation
//  3. Action            (child)              - Main command logic
//  4. PostRun           (child only)         - Command-specific cleanup
//  5. PersistentPostRun (child, then parent) - Global cleanup
//
// USAGE:
//
//	go run main.go --verbose              # Run root with all hooks visible
//	go run main.go sub --verbose          # Run subcommand, see hook cascade
//	go run main.go                        # Run without verbose (minimal output)
//
// KEY CONCEPTS:
//   - PersistentPreRun: Runs for command AND all subcommands (initialization)
//   - PreRun: Runs only for the executed command (validation)
//   - Action: The main command logic
//   - PostRun: Runs only for the executed command (cleanup)
//   - PersistentPostRun: Runs for command AND all subcommands (finalization)
//   - Hook inheritance: Child sees parent's persistent hooks
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nyxstack/cli"
)

var verbose bool

var rootCmd = cli.Root("lifecycle").
	Description("Demonstrates lifecycle hooks").
	Flag(&verbose, "verbose", "v", false, "Enable verbose output").
	PersistentPreRun(func(ctx context.Context, cmd *cli.Command) error {
		if verbose {
			fmt.Println("[Root PersistentPreRun] Initializing...")
		}
		return nil
	}).
	PreRun(func(ctx context.Context, cmd *cli.Command) error {
		if verbose {
			fmt.Println("[Root PreRun] Validating...")
		}
		return nil
	}).
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("[Root Action] Executing root command")
		return nil
	}).
	PostRun(func(ctx context.Context, cmd *cli.Command) error {
		if verbose {
			fmt.Println("[Root PostRun] Cleaning up...")
		}
		return nil
	}).
	PersistentPostRun(func(ctx context.Context, cmd *cli.Command) error {
		if verbose {
			fmt.Println("[Root PersistentPostRun] Finalizing...")
		}
		return nil
	})

var subCmd = cli.Cmd("sub").
	Description("Subcommand with hooks").
	PersistentPreRun(func(ctx context.Context, cmd *cli.Command) error {
		if verbose {
			fmt.Println("[Sub PersistentPreRun] Sub initializing...")
		}
		return nil
	}).
	PreRun(func(ctx context.Context, cmd *cli.Command) error {
		if verbose {
			fmt.Println("[Sub PreRun] Sub validating...")
		}
		return nil
	}).
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("[Sub Action] Executing subcommand")
		return nil
	}).
	PostRun(func(ctx context.Context, cmd *cli.Command) error {
		if verbose {
			fmt.Println("[Sub PostRun] Sub cleaning up...")
		}
		return nil
	}).
	PersistentPostRun(func(ctx context.Context, cmd *cli.Command) error {
		if verbose {
			fmt.Println("[Sub PersistentPostRun] Sub finalizing...")
		}
		return nil
	})

func init() {
	rootCmd.AddCommand(subCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
