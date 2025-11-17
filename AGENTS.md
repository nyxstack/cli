# AI Agent Guide for nyxstack/cli

Quick reference for AI agents to understand and use this CLI framework.

## What is nyxstack/cli?

A Go CLI framework with automatic flag inheritance and type-safe arguments. Simple API, powerful features.

## Core Concepts (30 seconds)

### 1. Basic Command
```go
cli.Root("myapp").
    Flag(&verbose, "verbose", "v", false, "Verbose output").
    Action(func(ctx context.Context, cmd *cli.Command) error {
        fmt.Println("Hello!")
        return nil
    }).
    Execute()
```

### 2. Flag Syntax
**CRITICAL:** Flags use `--flag=value` syntax (NOT `--flag value`)
```bash
myapp --port=8080        # ✅ Correct
myapp --port 8080        # ❌ Wrong
myapp --verbose=true     # ✅ Boolean with value
myapp -v                 # ✅ Boolean short form (no value needed)
```

### 3. Flag Inheritance
Parent flags automatically available on ALL child commands:
```go
root.Flag(&verbose, "verbose", "v", false, "Verbose")
child := cli.Cmd("child")
root.AddCommand(child)
// child automatically has --verbose flag!
```

### 4. Type-Safe Arguments
```go
cli.Cmd("deploy").
    Arg("environment", "Target environment", true).  // required
    Arg("replicas", "Number of replicas", true).     // required
    Action(func(ctx context.Context, cmd *cli.Command, env string, replicas int) error {
        // env is string, replicas is int (automatic conversion!)
        return nil
    })
```

## Quick API Reference

### Creating Commands
```go
root := cli.Root("myapp")           // Root command
cmd := cli.Cmd("subcommand")        // Subcommand
root.AddCommand(cmd)                // Add to hierarchy
```

### Flags
```go
// Basic flag
cmd.Flag(&port, "port", "p", 8080, "Port number")

// Required flag
cmd.FlagRequired(&apiKey, "api-key", "", "", "API key (required)")

// Hidden flag
cmd.FlagHidden(&debug, "debug", "", false, "Debug mode")

// Struct-based flags
type Config struct {
    Host string `cli:"host,h" default:"localhost" usage:"Server host"`
    Port int    `cli:"port,p" default:"8080" usage:"Server port"`
}
var config Config
cmd.Flags(&config)
```

### Arguments
```go
// Required argument
cmd.Arg("name", "User name", true)

// Optional argument
cmd.Arg("age", "User age", false)

// Variadic arguments
cmd.Arg("files", "Files to process", false)  // In action: files ...string
```

### Lifecycle Hooks
```go
cmd.
    PersistentPreRun(func(ctx context.Context, cmd *cli.Command) error {
        // Runs for this command + all children, before action
        return nil
    }).
    PreRun(func(ctx context.Context, cmd *cli.Command) error {
        // Runs only for this command, before action (good for validation)
        return nil
    }).
    Action(func(ctx context.Context, cmd *cli.Command) error {
        // Main command logic
        return nil
    }).
    PostRun(func(ctx context.Context, cmd *cli.Command) error {
        // Runs only for this command, after action
        return nil
    }).
    PersistentPostRun(func(ctx context.Context, cmd *cli.Command) error {
        // Runs for this command + all children, after action (good for cleanup)
        return nil
    })
```

**Execution order:** PersistentPreRun → PreRun → Action → PostRun → PersistentPostRun

### Context Support
```go
// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
app.ExecuteContext(ctx)

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
app.ExecuteContext(ctx)

// Pass values
ctx = context.WithValue(ctx, "key", value)
```

### Execution
```go
app.Execute()                              // Uses os.Args
app.ExecuteContext(ctx)                    // With context
app.ExecuteWithArgs([]string{"arg1"})      // Custom args (for testing)
```

## Common Patterns

### Complete CLI Application
```go
package main

import (
    "context"
    "fmt"
    "os"
    "github.com/nyxstack/cli"
)

func main() {
    var verbose bool
    var config string
    
    root := cli.Root("myapp").
        Description("My CLI application").
        Flag(&verbose, "verbose", "v", false, "Verbose output").
        Flag(&config, "config", "c", "config.yaml", "Config file")
    
    deploy := cli.Cmd("deploy").
        Description("Deploy the application").
        Arg("environment", "Target environment", true).
        PreRun(func(ctx context.Context, cmd *cli.Command) error {
            // Validate
            if !isValidEnv(environment) {
                return fmt.Errorf("invalid environment: %s", environment)
            }
            return nil
        }).
        Action(func(ctx context.Context, cmd *cli.Command, env string) error {
            if verbose {
                fmt.Printf("Deploying to %s...\n", env)
            }
            return doDeploy(env)
        })
    
    root.AddCommand(deploy)
    
    if err := root.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

### With Shell Completion
```go
root := cli.Root("myapp")
// ... add commands ...
cli.AddCompletion(root)  // Adds bash, zsh, fish, powershell completion
```

### Testing
```go
func TestCommand(t *testing.T) {
    var result string
    
    cmd := cli.Root("test").
        Action(func(ctx context.Context, cmd *cli.Command) error {
            result = "executed"
            return nil
        })
    
    err := cmd.ExecuteWithArgs([]string{})
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if result != "executed" {
        t.Error("command not executed")
    }
}
```

## Key Differences from Other Frameworks

| Feature | nyxstack/cli | Cobra | urfave/cli |
|---------|--------------|-------|------------|
| Flag inheritance | Automatic | Manual (Persistent) | Manual |
| Flag syntax | `--flag=value` | `--flag value` | `--flag value` |
| Type-safe args | Yes | No | No |
| Context support | Built-in | Manual | Manual |
| API style | Fluent/chainable | Struct-based | Struct-based |

## Error Handling

### Custom Error Types
```go
// CommandNotFoundError - command not found
// ArgumentError - missing/invalid argument
// FlagError - missing/invalid flag

if err := app.Execute(); err != nil {
    if cmdErr, ok := err.(*cli.CommandNotFoundError); ok {
        fmt.Printf("Unknown command: %s\n", cmdErr.Name)
        os.Exit(127)
    }
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(1)
}
```

### Return Errors from Actions
```go
cmd.Action(func(ctx context.Context, cmd *cli.Command) error {
    if err := doWork(); err != nil {
        return fmt.Errorf("work failed: %w", err)  // Wrap errors
    }
    return nil
})
```

## Important Gotchas

### 1. Flag Syntax is Strict
```go
// ❌ This will NOT work:
myapp --port 8080

// ✅ Must use equals:
myapp --port=8080
```

### 2. Flag Variables Must Be Pointers
```go
var port int
cmd.Flag(&port, "port", "p", 8080, "Port")  // Note the &port
```

### 3. Arguments Match Action Parameters
```go
cmd.
    Arg("name", "Name", true).
    Arg("count", "Count", true).
    Action(func(ctx context.Context, cmd *cli.Command, name string, count int) error {
        // Arguments must match in order and type
        return nil
    })
```

### 4. FlagSet Uses Pointers
When working with FlagSet directly, remember flags are stored as `[]*Flag` pointers.

### 5. Required Flags Checked Before Action
Required flags validated automatically before Action runs. No need to check in Action.

## File Structure

```
project/
├── main.go              # Entry point
├── go.mod              
└── cmd/                # Optional: command implementations
    ├── deploy.go
    ├── rollback.go
    └── server.go
```

Example main.go:
```go
package main

import (
    "github.com/nyxstack/cli"
    "myapp/cmd"
)

func main() {
    root := cli.Root("myapp")
    root.AddCommand(cmd.DeployCommand())
    root.AddCommand(cmd.RollbackCommand())
    root.AddCommand(cmd.ServerCommand())
    root.Execute()
}
```

## Supported Types

### Flags and Arguments
- `string`
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`
- `bool`
- `time.Duration`
- `[]string`, `[]int`, etc. (arrays via repeated flags)

## Testing Quick Reference

```go
// Test with custom args
err := cmd.ExecuteWithArgs([]string{"--flag=value", "arg1", "arg2"})

// Test with context
ctx := context.WithValue(context.Background(), "key", "value")
err := cmd.ExecuteContext(ctx)

// Test error types
if _, ok := err.(*cli.ArgumentError); !ok {
    t.Error("expected ArgumentError")
}
```

## Dependencies

- **Go 1.20+** required
- **github.com/nyxstack/color** - only dependency (for terminal colors)

## Installation

```bash
go get github.com/nyxstack/cli
```

## Documentation

- **Quick Start:** [docs/quick-start.md](docs/quick-start.md)
- **API Reference:** [docs/api-reference.md](docs/api-reference.md)
- **Examples:** [examples/](examples/)
- **Full Docs:** [docs/](docs/)

## Coverage

- **85.9% test coverage**
- **79 test functions**
- **All core features tested**

## When to Use This Framework

**Good for:**
- CLI tools with subcommands (git-style)
- Apps needing flag inheritance
- Tools with complex command hierarchies
- Type-safe argument handling
- Context-aware applications

**Not ideal for:**
- Simple single-command tools (overkill)
- Need POSIX-compliant `--flag value` syntax
- Maximum flexibility over simplicity

## Quick Decision Tree

```
Need CLI framework?
├─ Single command only → Use flag package
├─ Multiple subcommands → Use nyxstack/cli
│  ├─ Need automatic flag inheritance? → nyxstack/cli ✅
│  ├─ Need type-safe arguments? → nyxstack/cli ✅
│  └─ Prefer struct-based config? → Consider Cobra/urfave
└─ Need POSIX compliance → Use getopt library
```

## Example: Full-Featured Command

```go
var (
    verbose bool
    timeout time.Duration
    env     string
    dryRun  bool
)

cmd := cli.Root("deploy").
    Description("Deploy application to target environment").
    
    // Flags
    Flag(&verbose, "verbose", "v", false, "Enable verbose logging").
    Flag(&timeout, "timeout", "t", 5*time.Minute, "Deployment timeout").
    FlagRequired(&env, "env", "e", "", "Target environment (prod/staging/dev)").
    Flag(&dryRun, "dry-run", "", false, "Run without making changes").
    
    // Arguments
    Arg("service", "Service name to deploy", true).
    Arg("version", "Version tag to deploy", true).
    
    // Hooks
    PersistentPreRun(func(ctx context.Context, cmd *cli.Command) error {
        // Initialize (runs for this and all children)
        return initConfig()
    }).
    PreRun(func(ctx context.Context, cmd *cli.Command) error {
        // Validate (runs only for this command)
        if env != "prod" && env != "staging" && env != "dev" {
            return fmt.Errorf("invalid environment: %s", env)
        }
        return nil
    }).
    Action(func(ctx context.Context, cmd *cli.Command, service, version string) error {
        // Main logic
        if verbose {
            fmt.Printf("Deploying %s:%s to %s\n", service, version, env)
        }
        
        if dryRun {
            fmt.Println("Dry run - no changes made")
            return nil
        }
        
        ctx, cancel := context.WithTimeout(ctx, timeout)
        defer cancel()
        
        return deploy(ctx, service, version, env)
    }).
    PostRun(func(ctx context.Context, cmd *cli.Command) error {
        // Cleanup (runs only for this command)
        return sendNotification()
    })
```

## That's It!

You now know enough to build CLI applications with nyxstack/cli. Check [examples/](examples/) for more patterns.
