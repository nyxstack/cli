// Input Validation Example
//
// WHAT: Demonstrates input validation using PreRun hook before action execution.
//
// WHY: PreRun is the ideal place for validation because:
//   - Runs after flag parsing but before action
//   - Can return error to stop execution
//   - Keeps validation logic separate from business logic
//   - Provides early failure with clear error messages
//
// VALIDATION PATTERNS:
//   - Enum validation: Check value against allowed list
//   - Range validation: Ensure numeric values are within bounds
//   - Format validation: Verify string patterns (emails, URLs, etc.)
//   - Dependency validation: Check flag combinations
//
// USAGE:
//
//	go run main.go --format=json --port=8080          # Valid inputs
//	go run main.go --format=xml                       # Error: invalid format
//	go run main.go --port=99999                       # Error: port out of range
//	go run main.go --format=yaml -p=443               # Valid with short flag
//
// KEY CONCEPTS:
//   - PreRun hook: Perfect for validation logic
//   - Early failure: Stop before Action if validation fails
//   - Clear errors: Return descriptive error messages
//   - Separation of concerns: Validation vs business logic
package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/nyxstack/cli"
)

var (
	format string
	port   int
)

var rootCmd = cli.Root("validate").
	Description("Demonstrates input validation using PreRun").
	Flag(&format, "format", "f", "json", "Output format (json, yaml, text)").
	Flag(&port, "port", "p", 8080, "Server port").
	PreRun(func(ctx context.Context, cmd *cli.Command) error {
		// Validate format
		validFormats := []string{"json", "yaml", "text"}
		valid := false
		for _, f := range validFormats {
			if f == format {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid format '%s', must be one of: %s",
				format, strings.Join(validFormats, ", "))
		}

		// Validate port
		if port < 1 || port > 65535 {
			return fmt.Errorf("port must be between 1 and 65535, got %d", port)
		}

		return nil
	}).
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Printf("Format: %s\n", format)
		fmt.Printf("Port: %d\n", port)
		fmt.Println("Validation passed!")
		return nil
	})

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
