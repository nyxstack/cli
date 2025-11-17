// Hidden Commands and Flags Example
//
// WHAT: Demonstrates hiding commands and flags from help and completion.
//
// WHY: Hidden features are useful for:
//   - Debug commands for developers (not for end users)
//   - Internal/experimental commands
//   - Deprecated commands (still work but not advertised)
//   - Sensitive flags (API keys, tokens)
//   - Administrative features
//
// BEHAVIOR:
//   - Hidden commands: Don't appear in --help or shell completions
//   - Still executable: If you know the name, you can run it
//   - Hidden flags: Work but not shown in help text
//   - Use for internal tooling
//
// USAGE:
//
//	go run main.go --help                   # 'secret' command NOT shown
//	go run main.go public                   # Public command works
//	go run main.go secret                   # Hidden command still works!
//	go run main.go --debug                  # Hidden flag works
//	go run main.go --api-key=abc123         # Hidden flag for sensitive data
//
// KEY CONCEPTS:
//   - .Hidden() method: Mark command as hidden
//   - Hidden flags: Marked internally (future: .Hidden() on flags)
//   - Still functional: Hidden doesn't mean disabled
//   - Not in completions: Won't appear in shell autocomplete
//   - Use cases: Debug tools, deprecated features, admin commands
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nyxstack/cli"
)

var (
	verbose bool
	debug   bool
	apiKey  string
)

var rootCmd = cli.Root("hidden-flags").
	Description("Demonstrates hidden flags and commands").
	Flag(&verbose, "verbose", "v", false, "Enable verbose output").
	Flag(&debug, "debug", "d", false, "Enable debug mode (hidden)").
	Flag(&apiKey, "api-key", "", "", "API key for authentication (hidden)").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("Application running")
		if verbose {
			fmt.Println("  Verbose: enabled")
		}
		if debug {
			fmt.Println("  Debug: enabled (hidden flag)")
		}
		if apiKey != "" {
			fmt.Printf("  API Key: %s (hidden flag)\n", apiKey)
		}
		return nil
	})

// Public command - appears in help and completions
var publicCmd = cli.Cmd("public").
	Description("A public command visible in help").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("Public command executed")
		return nil
	})

// Hidden command - doesn't appear in help or completions
// Still executable if you know the name
var secretCmd = cli.Cmd("secret").
	Description("A hidden command for internal use").
	Hidden().
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("Secret command executed")
		if debug {
			fmt.Println("Debug mode is on")
		}
		return nil
	})

func init() {
	rootCmd.AddCommand(publicCmd)
	rootCmd.AddCommand(secretCmd)

	// Note: These flags would be marked as hidden in the implementation
	// The design shows .Hidden() method for flags (future enhancement)
	// For now, they work but will appear in --help
	// Implementation will add: flag.Hidden = true
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
