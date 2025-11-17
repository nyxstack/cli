# nyxstack/cli

A simple, powerful CLI framework for Go with automatic flag inheritance and type-safe arguments.

[![Go Version](https://img.shields.io/badge/go-1.20+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Test Coverage](https://img.shields.io/badge/coverage-85.9%25-brightgreen.svg)](TESTING.md)

## Features

- **Automatic Flag Inheritance** - Parent flags available to all child commands
- **Type-Safe Arguments** - Automatic type conversion for function parameters
- **Lifecycle Hooks** - PreRun, PostRun, PersistentPreRun, PersistentPostRun
- **Context Support** - Built-in context.Context for cancellation and timeouts
- **Shell Completion** - Bash, Zsh, Fish, and PowerShell
- **Fluent API** - Chainable interface for building CLI apps
- **Zero Dependencies** - Only depends on `github.com/nyxstack/color`

## Installation

```bash
go get github.com/nyxstack/cli
```

## Documentation

📖 **[Full Documentation](docs/)** - Comprehensive guides and API reference

- [Quick Start](docs/quick-start.md) - Get started in 5 minutes
- [Commands](docs/commands.md) - Creating and organizing commands
- [Flags & Arguments](docs/flags-and-arguments.md) - Working with flags and arguments
- [Lifecycle Hooks](docs/lifecycle.md) - PreRun, PostRun, and execution flow
- [API Reference](docs/api-reference.md) - Complete API documentation

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "github.com/nyxstack/cli"
)

func main() {
    var verbose bool
    
    app := cli.Root("myapp").
        Description("My awesome CLI application").
        Flag(&verbose, "verbose", "v", false, "Enable verbose output").
        Action(func(ctx context.Context, cmd *cli.Command) error {
            if verbose {
                fmt.Println("Verbose mode enabled")
            }
            fmt.Println("Hello from myapp!")
            return nil
        })
    
    if err := app.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, "Error:", err)
        os.Exit(1)
    }
}
```

**Usage:**
```bash
myapp                      # Hello from myapp!
myapp --verbose=true       # Verbose mode enabled + Hello from myapp!
myapp -v                   # Same as above
myapp --help               # Shows automatic help
```

See the **[Quick Start Guide](docs/quick-start.md)** for more examples.

## Why nyxstack/cli?

### Automatic Flag Inheritance
```go
var verbose bool
root.Flag(&verbose, "verbose", "v", false, "Verbose output")

deploy := cli.Cmd("deploy")
root.AddCommand(deploy)

// --verbose automatically available on deploy command!
// Works: myapp --verbose=true deploy
// Works: myapp deploy --verbose=true
```

### Type-Safe Arguments
```go
cli.Cmd("deploy").
    Arg("environment", "Target environment", true).
    Arg("replicas", "Number of replicas", true).
    Action(func(ctx context.Context, cmd *cli.Command, env string, replicas int) error {
        // env is string, replicas is int - automatic conversion!
        fmt.Printf("Deploying to %s with %d replicas\n", env, replicas)
        return nil
    })
```

### Context Support
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

app.ExecuteContext(ctx)  // Automatic timeout handling
```

## Examples

**Comprehensive examples in [examples/](examples/) directory:**

- [basic](examples/basic/) - Simple commands and flags
- [subcommands](examples/subcommands/) - Command hierarchies
- [flags-struct](examples/flags-struct/) - Struct-based flags
- [lifecycle](examples/lifecycle/) - Hooks and execution flow
- [validation](examples/validation/) - Input validation
- [inheritance](examples/inheritance/) - Flag inheritance
- And 8 more...

## Contributing

Contributions welcome! Please ensure:
1. All tests pass: `go test -v`
2. Coverage maintained: `go test -cover`
3. Code formatted: `go fmt ./...`

## License

MIT License - see [LICENSE](LICENSE) file.
