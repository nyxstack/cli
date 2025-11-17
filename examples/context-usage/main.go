// Context Usage Patterns Example
//
// WHAT: Demonstrates proper context usage for timeouts, cancellation, and values.
//
// WHY: Context is essential for:
//   - Cancellation: Stop long-running operations gracefully
//   - Timeouts: Enforce time limits on operations
//   - Request-scoped values: Pass data without global variables
//   - Resource cleanup: Defer cancellation to prevent leaks
//   - Propagation: Pass cancellation signals through call stack
//
// CONTEXT PATTERNS:
//  1. WithTimeout:   Set operation deadline
//  2. WithCancel:    Programmatic cancellation
//  3. WithValue:     Request-scoped data (user ID, request ID, etc.)
//  4. Combined:      Mix timeout + values for real-world scenarios
//
// USAGE:
//
//	go run main.go timeout --timeout 3s          # Times out after 3s
//	go run main.go cancel                        # Cancels after 2s
//	go run main.go values --user admin           # Context with values
//	go run main.go combined --timeout 5s         # Timeout + values
//
// KEY CONCEPTS:
//   - Always defer cancel(): Prevents context leaks
//   - Check ctx.Done(): Respect cancellation in loops
//   - Use select: Handle cancellation + work concurrently
//   - Context values: For request-scoped data only (not for optional parameters)
//   - Propagation: Pass ctx to functions that need cancellation
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nyxstack/cli"
)

var (
	timeout  time.Duration
	userID   string
	logLevel string
)

var rootCmd = cli.Root("context-demo").
	Description("Demonstrates context usage patterns").
	Flag(&timeout, "timeout", "t", 30*time.Second, "Operation timeout").
	Flag(&userID, "user", "u", "anonymous", "User ID").
	Flag(&logLevel, "log-level", "l", "info", "Log level")

// Context with timeout
var timeoutCmd = cli.Cmd("timeout").
	Description("Demonstrates context timeout").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		// Create timeout context
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		fmt.Printf("Starting operation with %s timeout...\n", timeout)

		// Simulate long-running operation
		select {
		case <-time.After(5 * time.Second):
			fmt.Println("Operation completed successfully")
			return nil
		case <-ctx.Done():
			return fmt.Errorf("operation timed out: %v", ctx.Err())
		}
	})

// Context with cancellation
var cancelCmd = cli.Cmd("cancel").
	Description("Demonstrates context cancellation").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		ctx, cancel := context.WithCancel(ctx)

		// Simulate cancellation after 2 seconds
		go func() {
			time.Sleep(2 * time.Second)
			fmt.Println("\nCancelling operation...")
			cancel()
		}()

		fmt.Println("Starting operation (will be cancelled)...")

		// Simulate work that respects cancellation
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				return fmt.Errorf("operation cancelled: %v", ctx.Err())
			case <-time.After(500 * time.Millisecond):
				fmt.Printf("  Step %d/10\n", i+1)
			}
		}

		return nil
	})

// Context with values (request-scoped data)
var valuesCmd = cli.Cmd("values").
	Description("Demonstrates context values for request-scoped data").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		// Add request-scoped values to context
		ctx = context.WithValue(ctx, "userID", userID)
		ctx = context.WithValue(ctx, "requestID", "req-12345")
		ctx = context.WithValue(ctx, "logLevel", logLevel)

		// Pass context to helper functions
		processRequest(ctx)
		logActivity(ctx)

		return nil
	})

// Helper function that uses context values
func processRequest(ctx context.Context) {
	userID := ctx.Value("userID").(string)
	requestID := ctx.Value("requestID").(string)

	fmt.Printf("Processing request:\n")
	fmt.Printf("  Request ID: %s\n", requestID)
	fmt.Printf("  User ID: %s\n", userID)
}

// Another helper that uses context values
func logActivity(ctx context.Context) {
	userID := ctx.Value("userID").(string)
	logLevel := ctx.Value("logLevel").(string)

	fmt.Printf("\nLogging activity:\n")
	fmt.Printf("  User: %s\n", userID)
	fmt.Printf("  Log Level: %s\n", logLevel)
}

// Combining timeout with values
var combinedCmd = cli.Cmd("combined").
	Description("Combines context timeout and values").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		// Add values
		ctx = context.WithValue(ctx, "userID", userID)

		// Add timeout
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		fmt.Printf("User %s starting operation with %s timeout\n",
			ctx.Value("userID"), timeout)

		// Simulate work
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("Operation completed successfully")
			return nil
		case <-ctx.Done():
			return fmt.Errorf("operation failed: %v", ctx.Err())
		}
	})

func init() {
	rootCmd.AddCommand(timeoutCmd)
	rootCmd.AddCommand(cancelCmd)
	rootCmd.AddCommand(valuesCmd)
	rootCmd.AddCommand(combinedCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
