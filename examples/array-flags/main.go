// Array/Slice Flags Example
//
// WHAT: Demonstrates flags that accept multiple values (arrays/slices).
//
// WHY: Many scenarios require multiple values:
//   - Tags for deployment (--tag v1.0 --tag production --tag stable)
//   - Include/exclude patterns (--include *.go --exclude *_test.go)
//   - Multiple hosts (--host localhost --host 127.0.0.1)
//   - Environment variables (--env KEY1=val1 --env KEY2=val2)
//
// USAGE PATTERNS:
//  1. Repeated flags: --tag value1 --tag value2 --tag value3
//  2. Short flags:    -t value1 -t value2 -t value3
//  3. Mixed:          --tag value1 -t value2 --tag value3
//
// USAGE:
//
//	go run main.go deploy prod --tag v1.0 --tag stable              # Multiple tags
//	go run main.go search query -i "*.go" -i "*.md" -e "*_test.go"  # Include/exclude
//	go run main.go server --host 0.0.0.0 --host localhost           # Multiple hosts
//	go run main.go run cmd --env USER=admin --env DEBUG=true        # Environment vars
//
// KEY CONCEPTS:
//   - Slice type: Flag type is []string, []int, etc.
//   - Repeated flags: Each occurrence appends to slice
//   - Default empty: Default is usually []Type{}
//   - Check length: Always check len(slice) > 0 before using
//   - Common pattern: For lists of values that can grow
package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/nyxstack/cli"
)

var (
	tags     []string
	excludes []string
	includes []string
	verbose  bool
)

var rootCmd = cli.Root("array-flags").
	Description("Demonstrates array/slice flags for multiple values").
	Flag(&verbose, "verbose", "v", false, "Enable verbose output")

// Multiple tags
var deployCmd = cli.Cmd("deploy").
	Description("Deploy with multiple tags").
	Flag(&tags, "tag", "t", []string{}, "Add deployment tag (can be used multiple times)").
	Arg("environment", "Target environment", true).
	Action(func(ctx context.Context, cmd *cli.Command, env string) error {
		fmt.Printf("Deploying to %s\n", env)

		if len(tags) > 0 {
			fmt.Println("Tags:")
			for _, tag := range tags {
				fmt.Printf("  - %s\n", tag)
			}
		} else {
			fmt.Println("No tags specified")
		}

		return nil
	})

// Include/exclude patterns
var searchCmd = cli.Cmd("search").
	Description("Search with include/exclude patterns").
	Flag(&includes, "include", "i", []string{}, "Include pattern").
	Flag(&excludes, "exclude", "e", []string{}, "Exclude pattern").
	Arg("query", "Search query", true).
	Action(func(ctx context.Context, cmd *cli.Command, query string) error {
		fmt.Printf("Searching for: %s\n", query)

		if len(includes) > 0 {
			fmt.Println("\nInclude patterns:")
			for _, pattern := range includes {
				fmt.Printf("  + %s\n", pattern)
			}
		}

		if len(excludes) > 0 {
			fmt.Println("\nExclude patterns:")
			for _, pattern := range excludes {
				fmt.Printf("  - %s\n", pattern)
			}
		}

		if len(includes) == 0 && len(excludes) == 0 {
			fmt.Println("\nNo patterns specified, searching all")
		}

		return nil
	})

// Multiple values for configuration
var serverCmd = cli.Cmd("server").
	Description("Start server with multiple allowed hosts").
	Flag(&includes, "host", "h", []string{"localhost"}, "Allowed host (can specify multiple)").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("Starting server with allowed hosts:")
		for _, host := range includes {
			fmt.Printf("  - %s\n", host)
		}

		if verbose {
			fmt.Printf("\nTotal allowed hosts: %d\n", len(includes))
		}

		return nil
	})

// Environment variables (multiple key=value pairs)
var envVars []string

var runCmd = cli.Cmd("run").
	Description("Run command with environment variables").
	Flag(&envVars, "env", "e", []string{}, "Environment variable (KEY=VALUE)").
	Arg("command", "Command to run", true).
	Action(func(ctx context.Context, cmd *cli.Command, command string) error {
		fmt.Printf("Running: %s\n", command)

		if len(envVars) > 0 {
			fmt.Println("\nEnvironment variables:")
			for _, env := range envVars {
				parts := strings.SplitN(env, "=", 2)
				if len(parts) == 2 {
					fmt.Printf("  %s = %s\n", parts[0], parts[1])
				} else {
					fmt.Printf("  %s (invalid format)\n", env)
				}
			}
		}

		return nil
	})

func init() {
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(runCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
