// Struct-Based Flags Example
//
// WHAT: Demonstrates defining flags using struct tags instead of individual Flag() calls.
//
// WHY: For commands with many flags, struct-based definition is cleaner and more maintainable.
// Benefits:
//   - Group related flags in a struct
//   - Use struct tags for configuration
//   - Type-safe with compile-time checking
//   - Easy to pass configuration to other functions
//   - Supports complex types (time.Duration, custom types)
//
// USAGE:
//
// \tgo run main.go                                    # Use all defaults
// \tgo run main.go --host=0.0.0.0 --port=3000        # Custom host/port
// \tgo run main.go --tls --timeout=60s               # Enable TLS with custom timeout
// \tgo run main.go -h=localhost -p=8080 -l=debug     # Short flags
// \tgo run main.go --help                            # Show all flags
//
// KEY CONCEPTS:
//   - Struct tags: `cli:"name,short" default:"value" usage:"description"`
//   - Type conversion: Framework handles string -> int, bool, duration, etc.
//   - Defaults: Specified in tags, applied automatically
//   - Flags() method: Binds entire struct at once
//   - Configuration structs: Pass to other functions for clean API
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nyxstack/cli"
)

// ServerConfig demonstrates struct-based flag definition
type ServerConfig struct {
	Host     string        `cli:"host,h" default:"localhost" usage:"Server host"`
	Port     int           `cli:"port,p" default:"8080" usage:"Server port"`
	TLS      bool          `cli:"tls" default:"false" usage:"Enable TLS"`
	Timeout  time.Duration `cli:"timeout,t" default:"30s" usage:"Request timeout"`
	LogLevel string        `cli:"log-level,l" default:"info" usage:"Log level (debug, info, warn, error)"`
}

var config ServerConfig

var rootCmd = cli.Root("server").
	Description("Server with struct-based configuration").
	Flags(&config).
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("Server Configuration:")
		fmt.Printf("  Host: %s\n", config.Host)
		fmt.Printf("  Port: %d\n", config.Port)
		fmt.Printf("  TLS: %v\n", config.TLS)
		fmt.Printf("  Timeout: %s\n", config.Timeout)
		fmt.Printf("  Log Level: %s\n", config.LogLevel)

		fmt.Println("\nStarting server...")
		return nil
	})

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
