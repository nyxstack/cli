// Basic CLI Example
//
// WHAT: Demonstrates the simplest possible CLI application with a root command and flags.
//
// WHY: Shows the fundamental building blocks of the CLI framework:
//   - Creating a root command with cli.Root()
//   - Adding flags with different types (bool, string)
//   - Defining an Action handler that executes the command logic
//   - Proper error handling with Execute() and os.Exit()
//
// USAGE:
//
//	go run main.go                          # Run without flags
//	go run main.go --verbose                # Enable verbose output
//	go run main.go -v --config=app.yaml     # Short flag + config file
//	go run main.go --help                   # Show help (automatic)
//
// KEY CONCEPTS:
//   - Root command: Entry point of the application
//   - Flags: Optional parameters with names, shortcuts, defaults, and descriptions
//   - Action: The main logic executed when the command runs
//   - Silent framework: No automatic output, you control everything
//   - Error handling: Execute() returns error, caller decides exit code
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nyxstack/cli"
)

var (
	verbose bool
	config  string
)

var rootCmd = cli.Root("basic").
	Description("A basic CLI example").
	Flag(&verbose, "verbose", "v", false, "Enable verbose output").
	Flag(&config, "config", "c", "", "Configuration file").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("Basic command executed!")
		if verbose {
			fmt.Println("Verbose mode enabled")
		}
		if config != "" {
			fmt.Printf("Using config: %s\n", config)
		}
		return nil
	})

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
