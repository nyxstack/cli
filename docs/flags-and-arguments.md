# Flags and Arguments

Learn how to use flags (options) and arguments (positional parameters) in your CLI applications.

## Flags

Flags are optional or required parameters that modify command behavior.

### Basic Flag Usage

```go
var verbose bool
var port int
var name string

cmd := cli.Root("myapp").
    Flag(&verbose, "verbose", "v", false, "Enable verbose output").
    Flag(&port, "port", "p", 8080, "Server port").
    Flag(&name, "name", "n", "default", "Service name")
```

### Flag Syntax

Flags **must** use the `--flag=value` format:

```bash
# Correct
myapp --port=8080
myapp -p=8080
myapp --verbose=true
myapp -v              # Boolean flags can omit =value

# Incorrect (will cause errors)
myapp --port 8080     # ❌ Space-separated not supported
```

**Why this syntax?** It eliminates ambiguity between flag values and arguments/subcommands. See [Flag Syntax](flag-syntax.md) for details.

### Flag Types

All basic Go types are supported:

```go
var (
    // Strings
    name    string
    
    // Integers
    port    int
    count   uint
    size    int64
    
    // Floating point
    ratio   float64
    percent float32
    
    // Boolean
    verbose bool
    
    // Duration
    timeout time.Duration
    
    // Arrays (repeated flags)
    tags    []string
)

cmd.Flag(&name, "name", "n", "default", "Name")
cmd.Flag(&port, "port", "p", 8080, "Port")
cmd.Flag(&ratio, "ratio", "r", 0.5, "Ratio")
cmd.Flag(&verbose, "verbose", "v", false, "Verbose")
cmd.Flag(&timeout, "timeout", "t", 30*time.Second, "Timeout")
cmd.Flag(&tags, "tag", "", nil, "Tags (repeatable)")
```

### Required Flags

Mark flags as required:

```go
var apiKey string
cmd.FlagRequired(&apiKey, "api-key", "", "", "API key (required)")

// Returns error if not provided
```

### Hidden Flags

Hide flags from help output (useful for deprecated or debug flags):

```go
var debugMode bool
cmd.FlagHidden(&debugMode, "debug", "", false, "Enable debug mode")
```

### Array Flags

Accept multiple values:

```go
var tags []string
cmd.Flag(&tags, "tag", "t", nil, "Add tag (repeatable)")
```

Usage:
```bash
myapp --tag=prod --tag=api --tag=v2
# tags = []string{"prod", "api", "v2"}
```

### Struct-Based Flags

Define flags using struct tags:

```go
type Config struct {
    Host    string        `cli:"host,h" default:"localhost" usage:"Server host"`
    Port    int           `cli:"port,p" default:"8080" usage:"Server port"`
    Verbose bool          `cli:"verbose,v" default:"false" usage:"Verbose output"`
    Timeout time.Duration `cli:"timeout,t" default:"30s" usage:"Timeout"`
}

var config Config
cmd.Flags(&config)
```

Struct tag format: `cli:"name,short"`

## Arguments

Arguments are positional parameters that must appear after flags.

### Basic Arguments

```go
cmd := cli.Cmd("deploy").
    Arg("environment", "Target environment", true).   // Required
    Arg("version", "Version to deploy", false).       // Optional
    Action(func(ctx context.Context, cmd *cli.Command, env, version string) error {
        if version == "" {
            version = "latest"
        }
        fmt.Printf("Deploying %s to %s\n", version, env)
        return nil
    })
```

Usage:
```bash
myapp deploy production        # env=production, version=""
myapp deploy production v2.1   # env=production, version=v2.1
myapp deploy                   # Error: missing required argument
```

### Argument Types

Arguments are automatically converted to function parameter types:

```go
cmd.Arg("count", "Number of items", true).
    Arg("ratio", "Success ratio", true).
    Arg("enabled", "Enable feature", true).
    Action(func(ctx context.Context, cmd *cli.Command, count int, ratio float64, enabled bool) error {
        // count is int
        // ratio is float64
        // enabled is bool
        return nil
    })
```

Supported types:
- `string`
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`
- `bool`
- `time.Duration`

### Variadic Arguments

Accept variable number of arguments:

```go
cmd.Arg("files", "Files to process", true).
    Action(func(ctx context.Context, cmd *cli.Command, files ...string) error {
        for _, file := range files {
            fmt.Println("Processing", file)
        }
        return nil
    })
```

Usage:
```bash
myapp process file1.txt file2.txt file3.txt
# files = []string{"file1.txt", "file2.txt", "file3.txt"}
```

### Mixed Arguments

Combine required, optional, and variadic:

```go
cmd.Arg("command", "Command to run", true).      // Required
    Arg("target", "Target server", true).        // Required
    Arg("args", "Command arguments", false).     // Optional variadic
    Action(func(ctx context.Context, cmd *cli.Command, command, target string, args ...string) error {
        fmt.Printf("Running %s on %s with args: %v\n", command, target, args)
        return nil
    })
```

Usage:
```bash
myapp exec restart api-server
# command=restart, target=api-server, args=[]

myapp exec deploy api-server --force --env=prod
# command=deploy, target=api-server, args=[--force, --env=prod]
```

## Flag Inheritance

Child commands automatically inherit parent flags:

```go
var verbose bool
root := cli.Root("myapp").
    Flag(&verbose, "verbose", "v", false, "Verbose output")

deploy := cli.Cmd("deploy").
    Action(func(ctx context.Context, cmd *cli.Command) error {
        if verbose {  // Inherited from parent!
            fmt.Println("Verbose mode")
        }
        return nil
    })

root.AddCommand(deploy)
```

Usage:
```bash
myapp --verbose=true deploy
myapp deploy --verbose=true
# Both work! --verbose is available everywhere
```

See [Flag Inheritance](flag-inheritance.md) for details.

## Flag Shadowing

Child commands can override parent flags:

```go
var globalPort int
root.Flag(&globalPort, "port", "p", 8080, "Default port")

var serverPort int
server := cli.Cmd("server")
server.Flag(&serverPort, "port", "p", 3000, "Server-specific port")
```

The child flag takes precedence when the subcommand is used.

## Validation

Validate flags in PreRun hooks:

```go
cmd.PreRun(func(ctx context.Context, cmd *cli.Command) error {
    if port < 1 || port > 65535 {
        return fmt.Errorf("port must be between 1 and 65535")
    }
    if timeout < 0 {
        return fmt.Errorf("timeout cannot be negative")
    }
    return nil
})
```

## Best Practices

### Flag Naming

- Use lowercase with hyphens: `--api-key`
- Provide short forms for common flags: `-v`, `-p`, `-h`
- Keep names descriptive: `--timeout` not `--t`
- Use consistent naming across commands

### Flag Documentation

- Write clear usage text: "Enable verbose logging"
- Document default values: "Port number (default: 8080)"
- Explain when to use: "API key from dashboard (required for production)"

### Flag Organization

```go
// Good: Group related flags
cmd.
    // Server configuration
    Flag(&host, "host", "h", "localhost", "Server host").
    Flag(&port, "port", "p", 8080, "Server port").
    
    // Logging
    Flag(&verbose, "verbose", "v", false, "Verbose logging").
    Flag(&logFile, "log-file", "l", "", "Log file path").
    
    // Authentication
    FlagRequired(&apiKey, "api-key", "", "", "API key").
    Flag(&apiSecret, "api-secret", "", "", "API secret")
```

### Argument Order

- Required arguments first
- Optional arguments next
- Variadic argument last

```go
cmd.
    Arg("required1", "First required arg", true).
    Arg("required2", "Second required arg", true).
    Arg("optional1", "Optional arg", false).
    Arg("variadic", "Variadic args", false)  // Must be last
```

## Examples

See the [examples/](../examples/) directory for:
- [basic](../examples/basic/) - Simple flags
- [flags-struct](../examples/flags-struct/) - Struct-based flags
- [array-flags](../examples/array-flags/) - Multiple values
- [argument-types](../examples/argument-types/) - Type conversion
- [validation](../examples/validation/) - Input validation
- [inheritance](../examples/inheritance/) - Flag inheritance
- [flag-shadowing](../examples/flag-shadowing/) - Override parent flags

## Next Steps

- Learn about [Flag Inheritance](flag-inheritance.md)
- Understand [Lifecycle Hooks](lifecycle.md) for validation
- Read about [Flag Syntax](flag-syntax.md) design decisions
