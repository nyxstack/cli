// Error Handling in Lifecycle Hooks Example
//
// WHAT: Demonstrates how errors in lifecycle hooks affect execution flow.
//
// WHY: Understanding error propagation is critical for:
//   - Preventing invalid operations
//   - Ensuring proper cleanup even on failure
//   - Providing clear error messages to users
//   - Controlling execution flow
//
// ERROR FLOW RULES:
//  1. PersistentPreRun error → Stop immediately, run Post hooks
//  2. PreRun error → Stop immediately, run Post hooks
//  3. Action error → Run Post hooks, then return error
//  4. PostRun error → Run PersistentPostRun, then return error
//  5. PostRun ALWAYS runs (even if Action fails)
//  6. PersistentPostRun ALWAYS runs (for cleanup)
//
// USAGE:
//
//	go run main.go validate test --fail-at prerun      # PreRun validation fails
//	go run main.go process --fail-at action            # Action fails, PostRun still runs
//	go run main.go cleanup --fail-at postrun           # PostRun cleanup fails
//	go run main.go chain --fail-at persistent-prerun   # Early failure, Post hooks run
//	go run main.go success                             # All hooks succeed
//
// KEY CONCEPTS:
//   - Early exit: PreRun errors stop before Action
//   - Cleanup guarantee: Post hooks run even on error
//   - Error propagation: First error is returned
//   - Validation pattern: Use PreRun for validation
//   - Cleanup pattern: Use PostRun for guaranteed cleanup
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nyxstack/cli"
)

var (
	skipValidation bool
	failAt         string
)

var rootCmd = cli.Root("errors").
	Description("Demonstrates error handling in lifecycle hooks").
	Flag(&skipValidation, "skip-validation", "", false, "Skip validation").
	Flag(&failAt, "fail-at", "", "", "Simulate failure at stage (prerun, action, postrun)")

// Example 1: PreRun validation errors
var validateCmd = cli.Cmd("validate").
	Description("Demonstrates PreRun validation errors").
	Arg("input", "Input value to validate", true).
	PreRun(func(ctx context.Context, cmd *cli.Command) error {
		if failAt == "prerun" {
			return fmt.Errorf("validation failed: simulated error in PreRun")
		}

		fmt.Println("[PreRun] Validating input...")
		return nil
	}).
	Action(func(ctx context.Context, cmd *cli.Command, input string) error {
		fmt.Printf("[Action] Processing: %s\n", input)
		return nil
	}).
	PostRun(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("[PostRun] Cleanup after validation")
		return nil
	})

// Example 2: Action errors stop execution
var processCmd = cli.Cmd("process").
	Description("Demonstrates Action errors").
	PreRun(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("[PreRun] Preparing to process")
		return nil
	}).
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("[Action] Starting processing...")

		if failAt == "action" {
			return fmt.Errorf("processing failed: simulated error in Action")
		}

		fmt.Println("[Action] Processing completed")
		return nil
	}).
	PostRun(func(ctx context.Context, cmd *cli.Command) error {
		// PostRun is ALWAYS called, even if Action fails
		fmt.Println("[PostRun] Cleanup (always runs)")
		return nil
	})

// Example 3: PostRun errors
var cleanupCmd = cli.Cmd("cleanup").
	Description("Demonstrates PostRun errors").
	PreRun(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("[PreRun] Setup")
		return nil
	}).
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("[Action] Work completed successfully")
		return nil
	}).
	PostRun(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("[PostRun] Starting cleanup...")

		if failAt == "postrun" {
			return fmt.Errorf("cleanup failed: simulated error in PostRun")
		}

		fmt.Println("[PostRun] Cleanup completed")
		return nil
	})

// Example 4: PersistentPreRun error propagation
var chainCmd = cli.Cmd("chain").
	Description("Demonstrates error propagation through hook chain").
	PersistentPreRun(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("[PersistentPreRun] Initializing...")

		if failAt == "persistent-prerun" {
			return fmt.Errorf("initialization failed: simulated error")
		}

		return nil
	}).
	PreRun(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("[PreRun] This won't run if PersistentPreRun fails")
		return nil
	}).
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("[Action] This won't run if PreRun fails")
		return nil
	}).
	PostRun(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("[PostRun] This always runs for cleanup")
		return nil
	}).
	PersistentPostRun(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("[PersistentPostRun] Final cleanup")
		return nil
	})

// Example 5: No error - successful execution
var successCmd = cli.Cmd("success").
	Description("Demonstrates successful execution with all hooks").
	PersistentPreRun(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("[PersistentPreRun] ✓")
		return nil
	}).
	PreRun(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("[PreRun] ✓")
		return nil
	}).
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("[Action] ✓")
		return nil
	}).
	PostRun(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("[PostRun] ✓")
		return nil
	}).
	PersistentPostRun(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("[PersistentPostRun] ✓")
		return nil
	})

func init() {
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(processCmd)
	rootCmd.AddCommand(cleanupCmd)
	rootCmd.AddCommand(chainCmd)
	rootCmd.AddCommand(successCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "\n❌ Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("\n✅ Command completed successfully")
}
