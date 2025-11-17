// Flag Inheritance Example
//
// WHAT: Demonstrates automatic flag inheritance from parent to child commands.
//
// WHY: Flag inheritance enables:
//   - Global flags (verbose, config) available to all subcommands
//   - No need to re-define common flags for each command
//   - Consistent behavior across command hierarchy
//   - Cleaner command definitions
//
// INHERITANCE RULES:
//   - Child commands automatically inherit ALL parent flags
//   - Inheritance is transitive (grandchild gets parent + grandparent flags)
//   - Each command adds its own specific flags
//   - No distinction between "local" and "persistent" flags (all inherit)
//
// USAGE:
//
//	go run main.go --verbose                              # Root with verbose
//	go run main.go server --verbose --config=app.yaml     # Server inherits verbose+config
//	go run main.go server start -v -c=app.yaml -d         # Start inherits all + detach
//	go run main.go server start --help                    # See all inherited flags
//
// KEY CONCEPTS:
//   - Automatic inheritance: All parent flags available to children
//   - Single FlagSet: Each command has one set, merged with ancestors
//   - Fluent API: Build command tree naturally
//   - Transitive: Deep hierarchies work correctly
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
	port    int
	host    string
	detach  bool
)

var rootCmd = cli.Root("myapp").
	Description("Demonstrates flag inheritance").
	Flag(&verbose, "verbose", "v", false, "Enable verbose output").
	Flag(&config, "config", "c", "", "Configuration file").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("Root command")
		if verbose {
			fmt.Println("  Verbose: true")
		}
		if config != "" {
			fmt.Printf("  Config: %s\n", config)
		}
		return nil
	})

var serverCmd = cli.Cmd("server").
	Description("Server command (inherits verbose and config)").
	Flag(&port, "port", "p", 8080, "Server port").
	Flag(&host, "host", "h", "localhost", "Server host").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("Server command")
		if verbose {
			fmt.Println("  Verbose: true (inherited from root)")
		}
		if config != "" {
			fmt.Printf("  Config: %s (inherited from root)\n", config)
		}
		fmt.Printf("  Host: %s\n", host)
		fmt.Printf("  Port: %d\n", port)
		return nil
	})

var startCmd = cli.Cmd("start").
	Description("Start server (inherits all parent flags)").
	Flag(&detach, "detach", "d", false, "Run in background").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("Start command")
		if verbose {
			fmt.Println("  Verbose: true (inherited from root)")
		}
		if config != "" {
			fmt.Printf("  Config: %s (inherited from root)\n", config)
		}
		fmt.Printf("  Host: %s (inherited from server)\n", host)
		fmt.Printf("  Port: %d (inherited from server)\n", port)
		if detach {
			fmt.Println("  Detach: true")
		}
		return nil
	})

func init() {
	serverCmd.AddCommand(startCmd)
	rootCmd.AddCommand(serverCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
