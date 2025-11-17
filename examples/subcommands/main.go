// Subcommands Example
//
// WHAT: Demonstrates building a nested command structure with multiple subcommands.
//
// WHY: Most CLI applications need hierarchical command organization (like git, docker, kubectl).
// This example shows:
//   - Root command with global flags (verbose)
//   - Multiple subcommands at different levels (deploy, server, server start, server stop)
//   - Arguments with required/optional types
//   - Flag inheritance from parent to child commands
//   - Organizing related commands under a parent (server management)
//
// USAGE:
//
//	go run main.go deploy dev v1.2.3          # Deploy with version
//	go run main.go deploy staging             # Deploy without optional version
//	go run main.go server start --port=3000   # Start server on custom port
//	go run main.go server start -d            # Start in detached mode
//	go run main.go server stop -v             # Stop with verbose output
//	go run main.go --help                     # Show all commands
//
// KEY CONCEPTS:
//   - Subcommands: Organize functionality into logical groups
//   - Arguments: Positional parameters (required vs optional)
//   - Flag inheritance: Child commands inherit parent flags (verbose from root)
//   - Command tree: server -> start/stop creates nested hierarchy
//   - Action signatures: Can accept (ctx, cmd) or (ctx, cmd, arg1, arg2, ...)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nyxstack/cli"
)

var (
	verbose bool
	port    int
	host    string
	detach  bool
)

var rootCmd = cli.Root("myapp").
	Description("Application with subcommands").
	Flag(&verbose, "verbose", "v", false, "Enable verbose output")

var deployCmd = cli.Cmd("deploy").
	Description("Deploy the application").
	Arg("environment", "Target environment (dev, staging, prod)", true).
	Arg("version", "Version to deploy", false).
	Action(func(ctx context.Context, cmd *cli.Command, env string, version string) error {
		if verbose {
			fmt.Println("Verbose: Deploying...")
		}
		fmt.Printf("Deploying to %s\n", env)
		if version != "" {
			fmt.Printf("Version: %s\n", version)
		}
		return nil
	})

var serverCmd = cli.Cmd("server").
	Description("Server management").
	Flag(&port, "port", "p", 8080, "Server port").
	Flag(&host, "host", "h", "localhost", "Server host")

var startCmd = cli.Cmd("start").
	Description("Start the server").
	Flag(&detach, "detach", "d", false, "Run in background").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		if verbose {
			fmt.Println("Verbose: Starting server...")
		}
		fmt.Printf("Starting server on %s:%d\n", host, port)
		if detach {
			fmt.Println("Running in background")
		}
		return nil
	})

var stopCmd = cli.Cmd("stop").
	Description("Stop the server").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		if verbose {
			fmt.Println("Verbose: Stopping server...")
		}
		fmt.Println("Stopping server...")
		return nil
	})

func init() {
	serverCmd.AddCommand(startCmd)
	serverCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(serverCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
