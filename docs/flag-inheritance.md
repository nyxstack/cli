# Flag Inheritance

One of the most powerful features of nyxstack/cli is automatic flag inheritance. Child commands automatically inherit all flags from their parent commands.

## How It Works

When you define a flag on a parent command, it's automatically available on all child commands:

```go
var verbose bool

root := cli.Root("myapp").
    Flag(&verbose, "verbose", "v", false, "Enable verbose output")

deploy := cli.Cmd("deploy").
    Action(func(ctx context.Context, cmd *cli.Command) error {
        if verbose {  // Automatically available!
            fmt.Println("Deploying with verbose output")
        }
        return nil
    })

root.AddCommand(deploy)
```

Usage - flag can be used anywhere:
```bash
myapp --verbose=true deploy      # Flag before subcommand
myapp deploy --verbose=true      # Flag after subcommand
myapp deploy -v                  # Short form works too
```

## Multi-Level Inheritance

Flags inherit through multiple levels:

```go
var verbose bool
var debug bool
var trace bool

root := cli.Root("myapp").
    Flag(&verbose, "verbose", "v", false, "Verbose output")

database := cli.Cmd("database").
    Flag(&debug, "debug", "d", false, "Debug mode")

migrate := cli.Cmd("migrate").
    Flag(&trace, "trace", "t", false, "Trace SQL")

database.AddCommand(migrate)
root.AddCommand(database)
```

The `migrate` command has access to ALL three flags:
```bash
myapp database migrate --verbose=true --debug=true --trace=true
```

Inheritance chain:
- `migrate` has: `trace` (own), `debug` (parent), `verbose` (grandparent)
- `database` has: `debug` (own), `verbose` (parent)
- `root` has: `verbose` (own)

## Flag Shadowing

Child commands can override parent flags with the same name:

```go
var globalPort int
root := cli.Root("myapp").
    Flag(&globalPort, "port", "p", 8080, "Default port")

var serverPort int
server := cli.Cmd("server").
    Flag(&serverPort, "port", "p", 3000, "Server-specific port").
    Action(func(ctx context.Context, cmd *cli.Command) error {
        fmt.Printf("Using port: %d\n", serverPort)
        return nil
    })

root.AddCommand(server)
```

Usage:
```bash
myapp --port=8080              # Uses globalPort (8080)
myapp server --port=3000       # Uses serverPort (3000, child overrides parent)
myapp server                   # Uses serverPort default (3000)
```

The child's flag takes precedence when the child command is executed.

## Isolation Between Siblings

Sibling commands don't share flags:

```go
var deployEnv string
deploy := cli.Cmd("deploy").
    Flag(&deployEnv, "env", "e", "production", "Deploy environment")

var rollbackVersion string
rollback := cli.Cmd("rollback").
    Flag(&rollbackVersion, "version", "v", "", "Rollback version")

root.AddCommand(deploy)
root.AddCommand(rollback)
```

- `deploy` has access to `--env` only
- `rollback` has access to `--version` only
- Neither can use the other's flags

```bash
myapp deploy --env=staging        # ✅ Works
myapp rollback --version=1.0      # ✅ Works
myapp deploy --version=1.0        # ❌ Error: unknown flag
myapp rollback --env=staging      # ❌ Error: unknown flag
```

## Hidden Flag Inheritance

Hidden flags also inherit:

```go
var debugMode bool
root.FlagHidden(&debugMode, "debug", "", false, "Debug mode")

cmd := cli.Cmd("test")
root.AddCommand(cmd)
```

The `test` command inherits `--debug`, but it won't appear in help text.

## Parent Isolation

Child flags don't affect parent commands:

```go
var childFlag string
child := cli.Cmd("child").
    Flag(&childFlag, "child-flag", "", "", "Child-specific flag")

root.AddCommand(child)
```

```bash
myapp --child-flag=value        # ❌ Error: unknown flag
myapp child --child-flag=value  # ✅ Works
```

## Use Cases

### Global Configuration

Define common flags once at the root:

```go
var (
    configFile string
    verbose    bool
    dryRun     bool
)

root := cli.Root("myapp").
    Flag(&configFile, "config", "c", "config.yaml", "Config file").
    Flag(&verbose, "verbose", "v", false, "Verbose output").
    Flag(&dryRun, "dry-run", "", false, "Dry run mode")

// All subcommands automatically have access to these flags
```

### Environment-Specific Flags

```go
var environment string

root.Flag(&environment, "env", "e", "production", "Environment")

// Every subcommand can use --env
deploy.Action(func(...) { /* use environment */ })
rollback.Action(func(...) { /* use environment */ })
test.Action(func(...) { /* use environment */ })
```

### Database Connection

```go
var (
    dbHost string
    dbPort int
    dbName string
)

root.
    Flag(&dbHost, "db-host", "", "localhost", "Database host").
    Flag(&dbPort, "db-port", "", 5432, "Database port").
    Flag(&dbName, "db-name", "", "myapp", "Database name")

// All database commands inherit connection settings
```

### Authentication

```go
var apiToken string

root.FlagRequired(&apiToken, "token", "t", "", "API token")

// All commands that need API access inherit the token
```

## Best Practices

### 1. Define Shared Flags on Parent

```go
// Good: Shared flags on root
root.Flag(&verbose, "verbose", "v", false, "Verbose")

// Bad: Repeating on each command
deploy.Flag(&verbose, "verbose", "v", false, "Verbose")
rollback.Flag(&verbose, "verbose", "v", false, "Verbose")
test.Flag(&verbose, "verbose", "v", false, "Verbose")
```

### 2. Use Shadowing for Defaults

```go
// Global default port
root.Flag(&port, "port", "p", 8080, "Port")

// Server-specific default
server.Flag(&serverPort, "port", "p", 3000, "Server port")

// API-specific default
api.Flag(&apiPort, "port", "p", 9000, "API port")
```

### 3. Group Related Flags

```go
// Server group
server := cli.Cmd("server").
    Flag(&serverHost, "host", "h", "localhost", "Server host").
    Flag(&serverPort, "port", "p", 8080, "Server port")

// Database group
database := cli.Cmd("database").
    Flag(&dbHost, "host", "h", "localhost", "DB host").
    Flag(&dbPort, "port", "p", 5432, "DB port")
```

Different `--host` and `--port` flags for different contexts.

### 4. Document Inheritance

```go
root := cli.Root("myapp").
    Description("MyApp - All subcommands inherit root flags").
    Flag(&verbose, "verbose", "v", false, "Verbose (available on all commands)")
```

## Testing Inheritance

Test that child commands have access to parent flags:

```go
func TestFlagInheritance(t *testing.T) {
    var verbose bool
    root := cli.Root("test").
        Flag(&verbose, "verbose", "v", false, "Verbose")
    
    var executed bool
    child := cli.Cmd("child").
        Action(func(ctx context.Context, cmd *cli.Command) error {
            executed = true
            if !verbose {
                t.Error("verbose should be true")
            }
            return nil
        })
    
    root.AddCommand(child)
    
    // Flag before subcommand
    err := root.ExecuteWithArgs([]string{"--verbose=true", "child"})
    if err != nil || !executed {
        t.Error("execution failed")
    }
}
```

## Comparison with Other Frameworks

| Framework | Inheritance | Syntax |
|-----------|-------------|--------|
| nyxstack/cli | ✅ Automatic | All parent flags available |
| Cobra | ⚠️ Manual | Must mark as "Persistent" |
| urfave/cli | ❌ No | Must repeat on each command |

## Examples

See the [examples/](../examples/) directory:
- [inheritance](../examples/inheritance/) - Flag inheritance demo
- [flag-shadowing](../examples/flag-shadowing/) - Override parent flags
- [subcommands](../examples/subcommands/) - Multi-level inheritance

## Next Steps

- Learn about [Flags and Arguments](flags-and-arguments.md)
- Understand [Commands](commands.md)
- Read the [API Reference](api-reference.md)
