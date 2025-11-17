// Flag Shadowing Example
//
// WHAT: Demonstrates child commands overriding (shadowing) parent flags.
//
// WHY: Sometimes child commands need different behavior for same flag name:
//   - Database operations need longer timeout than root default
//   - Server commands want different default format (json vs text)
//   - Subcommands need specialized versions of global flags
//
// SHADOWING RULES:
//   - Child flag with same name overrides parent flag
//   - Only the child's flag variable is used
//   - Other inherited flags work normally
//   - Shadowing creates command-specific behavior
//
// WHEN TO USE:
//   - Different defaults: Child needs different default value
//   - Different types: Child wants different semantics
//   - Specialized behavior: Child interprets flag differently
//
// USAGE:
//
//	go run main.go --timeout=30 --verbose                 # Root timeout=30
//	go run main.go database --timeout=60                  # DB timeout=60 (shadowed)
//	go run main.go database --verbose                     # Uses root's verbose (inherited)
//	go run main.go server --verbose --format=json         # Server shadows verbose+format
//	go run main.go query test --timeout=30                # Query uses root timeout (inherited)
//
// KEY CONCEPTS:
//   - Shadowing: Child flag masks parent flag with same name
//   - Separate variables: Each flag binds to its own variable
//   - Partial shadowing: Can shadow some flags, inherit others
//   - Help text: Shows child's flag description when shadowed
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nyxstack/cli"
)

var (
	// Root level flags
	timeout int
	verbose bool
	format  string

	// Database level flags
	dbTimeout int

	// Server level flags
	serverVerbose bool
	serverFormat  string
)

var rootCmd = cli.Root("shadow").
	Description("Demonstrates flag shadowing (child overrides parent)").
	Flag(&timeout, "timeout", "t", 30, "Default timeout in seconds").
	Flag(&verbose, "verbose", "v", false, "Enable verbose output").
	Flag(&format, "format", "f", "text", "Output format")

// Database command - shadows timeout flag
var dbCmd = cli.Cmd("database").
	Description("Database operations (has its own timeout)").
	Flag(&dbTimeout, "timeout", "t", 60, "Database operation timeout in seconds").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("Database command:")
		fmt.Printf("  Timeout: %d seconds (database-specific)\n", dbTimeout)
		fmt.Printf("  Verbose: %v (inherited from root)\n", verbose)
		fmt.Printf("  Format: %s (inherited from root)\n", format)

		fmt.Println("\nNote: The timeout flag is shadowed by database command")
		return nil
	})

// Server command - shadows verbose and format flags
var serverCmd = cli.Cmd("server").
	Description("Server management (has its own verbose and format)").
	Flag(&serverVerbose, "verbose", "v", false, "Server-specific verbose mode").
	Flag(&serverFormat, "format", "f", "json", "Server output format").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("Server command:")
		fmt.Printf("  Timeout: %d seconds (inherited from root)\n", timeout)
		fmt.Printf("  Verbose: %v (server-specific)\n", serverVerbose)
		fmt.Printf("  Format: %s (server-specific)\n", serverFormat)

		fmt.Println("\nNote: verbose and format flags are shadowed by server command")
		return nil
	})

// Query command - no shadowing, uses all inherited flags
var queryCmd = cli.Cmd("query").
	Description("Query operations (inherits all root flags)").
	Arg("query", "Query string", true).
	Action(func(ctx context.Context, cmd *cli.Command, query string) error {
		fmt.Printf("Query: %s\n", query)
		fmt.Printf("  Timeout: %d seconds (inherited)\n", timeout)
		fmt.Printf("  Verbose: %v (inherited)\n", verbose)
		fmt.Printf("  Format: %s (inherited)\n", format)

		fmt.Println("\nNote: No shadowing, using all inherited flags")
		return nil
	})

// Nested shadowing example
var startCmd = cli.Cmd("start").
	Description("Start database server (inherits db timeout)").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("Starting database server:")
		fmt.Printf("  Timeout: %d seconds (inherited from database command)\n", dbTimeout)
		fmt.Printf("  Verbose: %v (inherited from root)\n", verbose)

		fmt.Println("\nNote: Uses database's shadowed timeout, not root's")
		return nil
	})

func init() {
	dbCmd.AddCommand(startCmd)
	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(queryCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
