// Execute vs ExecuteContext Example
//
// WHAT: Demonstrates the three execution methods and their use cases.
//
// WHY: Different scenarios need different execution patterns:
//   - Execute():              Simple apps, uses os.Args, background context
//   - ExecuteContext(ctx):    Apps needing timeout/cancellation control
//   - ExecuteWithArgs(args):  Testing with custom arguments
//
// EXECUTION METHODS:
//  1. Execute()               → ExecuteContext(context.Background())
//  2. ExecuteContext(ctx)     → Full control over context
//  3. ExecuteWithArgs(args)   → For testing (custom args + background context)
//
// ALL METHODS:
//   - Return error (never call os.Exit internally)
//   - Caller decides exit code
//   - Silent framework (no automatic output)
//
// USAGE:
//
//	go run main.go basic                    # Simple execution
//	go run main.go fail                     # Returns error
//	go run main.go long --timeout 3s        # Times out after 3s
//	go run main.go server                   # Cancels after 5s
//
// KEY CONCEPTS:
//   - All return errors: Framework never calls os.Exit()
//   - Context propagation: Passed to all hooks and actions
//   - Timeout pattern: WithTimeout for operation limits
//   - Cancellation pattern: WithCancel for graceful shutdown
//   - Testing pattern: ExecuteWithArgs for unit tests
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nyxstack/cli"
)

var (
	timeout time.Duration
)

var rootCmd = cli.Root("execute-demo").
	Description("Demonstrates Execute vs ExecuteContext patterns").
	Flag(&timeout, "timeout", "t", 5*time.Second, "Command timeout")

// Example 1: Basic Execute() - uses os.Args, returns error
var basicCmd = cli.Cmd("basic").
	Description("Basic Execute pattern").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("Basic command running...")
		time.Sleep(1 * time.Second)
		fmt.Println("Done!")
		return nil
	})

// Example 2: Execute with error handling
var failCmd = cli.Cmd("fail").
	Description("Command that returns an error").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("This command will fail...")
		return fmt.Errorf("simulated failure")
	})

// Example 3: ExecuteContext with timeout
var longCmd = cli.Cmd("long").
	Description("Long-running command with context timeout").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("Starting long operation...")

		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				return fmt.Errorf("operation cancelled: %v", ctx.Err())
			case <-time.After(1 * time.Second):
				fmt.Printf("  Progress: %d/10\n", i+1)
			}
		}

		fmt.Println("Completed!")
		return nil
	})

// Example 4: ExecuteContext with cancellation
var serverCmd = cli.Cmd("server").
	Description("Simulates a server that can be cancelled").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("Server starting...")

		// Simulate server running
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("\nServer shutting down gracefully...")
				return ctx.Err()
			case <-ticker.C:
				fmt.Println("Server is running...")
			}
		}
	})

func init() {
	rootCmd.AddCommand(basicCmd)
	rootCmd.AddCommand(failCmd)
	rootCmd.AddCommand(longCmd)
	rootCmd.AddCommand(serverCmd)
}

func main() {
	// All examples use ExecuteContext to demonstrate proper patterns

	// For basic/fail commands: use background context
	if len(os.Args) > 1 && (os.Args[1] == "basic" || os.Args[1] == "fail") {
		if err := rootCmd.ExecuteContext(context.Background()); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// For long command: use timeout context
	if len(os.Args) > 1 && os.Args[1] == "long" {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		fmt.Printf("Running with %s timeout\n", timeout)
		if err := rootCmd.ExecuteContext(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// For server command: use cancellable context
	if len(os.Args) > 1 && os.Args[1] == "server" {
		ctx, cancel := context.WithCancel(context.Background())

		// Simulate graceful shutdown after 5 seconds
		go func() {
			time.Sleep(5 * time.Second)
			fmt.Println("\n[Main] Sending shutdown signal...")
			cancel()
		}()

		if err := rootCmd.ExecuteContext(ctx); err != nil {
			if err == context.Canceled {
				fmt.Println("Server stopped")
			} else {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		}
		return
	}

	// Default: show help or run command
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
