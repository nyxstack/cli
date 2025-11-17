// Argument Type Conversion Example
//
// WHAT: Demonstrates automatic type conversion for command arguments.
//
// WHY: Type safety and automatic conversion:
//   - Arguments arrive as strings from command line
//   - Framework converts to Go types automatically
//   - Type-safe function signatures
//   - Compile-time checking
//   - Cleaner action code (no manual parsing)
//
// SUPPORTED TYPES:
//   - string:        No conversion needed
//   - int:           Parse integer values
//   - bool:          Parse true/false, yes/no, 1/0
//   - float64:       Parse floating point numbers
//   - time.Duration: Parse duration strings (1s, 500ms, 2m30s)
//   - Custom types:  Implement conversion interface
//
// USAGE:
//
//	go run main.go greet Alice                    # String argument
//	go run main.go repeat "Hello" 3               # String + int
//	go run main.go toggle feature1 true           # String + bool
//	go run main.go calculate 3.14 2.5             # Two floats
//	go run main.go wait 2s                        # Duration
//	go run main.go config 8080 30s 5              # Mixed types with optional
//
// KEY CONCEPTS:
//   - Type conversion: Automatic based on action signature
//   - Required vs optional: Specified in Arg() definition
//   - Zero values: Optional args get type's zero value if not provided
//   - Action signature: func(ctx, cmd, arg1 type1, arg2 type2, ...)
//   - Type safety: Compile-time checking prevents errors
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nyxstack/cli"
)

var verbose bool

var rootCmd = cli.Root("arg-types").
	Description("Demonstrates argument type conversion").
	Flag(&verbose, "verbose", "v", false, "Enable verbose output")

// String arguments
var greetCmd = cli.Cmd("greet").
	Description("Greet someone by name").
	Arg("name", "Person's name", true).
	Action(func(ctx context.Context, cmd *cli.Command, name string) error {
		fmt.Printf("Hello, %s!\n", name)
		return nil
	})

// Integer arguments
var repeatCmd = cli.Cmd("repeat").
	Description("Repeat a message N times").
	Arg("message", "Message to repeat", true).
	Arg("count", "Number of times to repeat (integer)", true).
	Action(func(ctx context.Context, cmd *cli.Command, message string, count int) error {
		if verbose {
			fmt.Printf("Repeating '%s' %d times\n", message, count)
		}
		for i := 0; i < count; i++ {
			fmt.Println(message)
		}
		return nil
	})

// Boolean arguments
var toggleCmd = cli.Cmd("toggle").
	Description("Toggle a feature on/off").
	Arg("feature", "Feature name", true).
	Arg("enabled", "Enable or disable (true/false)", true).
	Action(func(ctx context.Context, cmd *cli.Command, feature string, enabled bool) error {
		status := "disabled"
		if enabled {
			status = "enabled"
		}
		fmt.Printf("Feature '%s' is now %s\n", feature, status)
		return nil
	})

// Float arguments
var calculateCmd = cli.Cmd("calculate").
	Description("Calculate with floating point numbers").
	Arg("value1", "First number", true).
	Arg("value2", "Second number", true).
	Action(func(ctx context.Context, cmd *cli.Command, val1 float64, val2 float64) error {
		sum := val1 + val2
		product := val1 * val2
		fmt.Printf("Sum: %.2f\n", sum)
		fmt.Printf("Product: %.2f\n", product)
		return nil
	})

// Duration arguments
var waitCmd = cli.Cmd("wait").
	Description("Wait for a specified duration").
	Arg("duration", "How long to wait (e.g., 1s, 500ms, 2m)", true).
	Action(func(ctx context.Context, cmd *cli.Command, duration time.Duration) error {
		fmt.Printf("Waiting for %s...\n", duration)
		time.Sleep(duration)
		fmt.Println("Done!")
		return nil
	})

// Mixed types with optional arguments
var configCmd = cli.Cmd("config").
	Description("Configure with mixed argument types").
	Arg("port", "Port number (int)", true).
	Arg("timeout", "Request timeout (duration)", true).
	Arg("retries", "Number of retries (int, optional)", false).
	Action(func(ctx context.Context, cmd *cli.Command, port int, timeout time.Duration, retries int) error {
		fmt.Println("Configuration:")
		fmt.Printf("  Port: %d\n", port)
		fmt.Printf("  Timeout: %s\n", timeout)
		if retries > 0 {
			fmt.Printf("  Retries: %d\n", retries)
		} else {
			fmt.Println("  Retries: default (0)")
		}
		return nil
	})

func init() {
	rootCmd.AddCommand(greetCmd)
	rootCmd.AddCommand(repeatCmd)
	rootCmd.AddCommand(toggleCmd)
	rootCmd.AddCommand(calculateCmd)
	rootCmd.AddCommand(waitCmd)
	rootCmd.AddCommand(configCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
