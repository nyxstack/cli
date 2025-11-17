# Commands

Commands are the building blocks of your CLI application. They represent actions users can perform.

## Creating Commands

### Root Command

Every CLI app starts with a root command:

```go
app := cli.Root("myapp").
    Description("My application description")
```

The root command is what users invoke directly: `myapp`

### Regular Commands

Add subcommands to create command hierarchies:

```go
deploy := cli.Cmd("deploy").
    Description("Deploy the application")

rollback := cli.Cmd("rollback").
    Description("Rollback to previous version")

app.AddCommand(deploy)
app.AddCommand(rollback)
```

Usage:
```bash
myapp deploy
myapp rollback
```

### Nested Commands

Create deep command hierarchies:

```go
database := cli.Cmd("database").
    Description("Database operations")

migrate := cli.Cmd("migrate").
    Description("Run migrations")

rollback := cli.Cmd("rollback").
    Description("Rollback migrations")

database.AddCommand(migrate)
database.AddCommand(rollback)
app.AddCommand(database)
```

Usage:
```bash
myapp database migrate
myapp database rollback
```

## Command Properties

### Description

Provide user-friendly help text:

```go
cmd := cli.Cmd("deploy").
    Description("Deploy application to target environment")
```

### Hidden Commands

Hide internal or debug commands from help:

```go
debug := cli.Cmd("debug").
    Description("Debug tools").
    Hidden()  // Won't appear in help

debug.Show()  // Make it visible again
```

### Command Actions

Define what happens when the command runs:

```go
cmd := cli.Cmd("deploy").
    Action(func(ctx context.Context, cmd *cli.Command) error {
        fmt.Println("Deploying...")
        return nil
    })
```

## Command Hierarchy

Commands can have multiple levels:

```go
// myapp server start --port=8080
root := cli.Root("myapp")
server := cli.Cmd("server")
start := cli.Cmd("start")

var port int
start.Flag(&port, "port", "p", 8080, "Server port").
    Action(func(ctx context.Context, cmd *cli.Command) error {
        fmt.Printf("Starting server on port %d\n", port)
        return nil
    })

server.AddCommand(start)
root.AddCommand(server)
```

## Command Methods

### Getters

```go
name := cmd.GetName()                 // Command name
desc := cmd.GetDescription()          // Command description
parent := cmd.GetParent()             // Parent command
children := cmd.GetCommands()         // Map of subcommands
args := cmd.GetArgs()                 // List of arguments
```

### State

```go
if cmd.IsHidden() {
    // Command is hidden
}
```

## Method Chaining

All methods return `*Command` for easy chaining:

```go
cmd := cli.Cmd("deploy").
    Description("Deploy the app").
    Flag(&env, "env", "e", "prod", "Environment").
    Arg("service", "Service name", true).
    PreRun(validateInput).
    Action(doDeploy).
    PostRun(cleanup)
```

## Execution Flow

When a command is executed, the following happens in order:

1. Parse flags and arguments
2. Run `PersistentPreRun` hooks (parent → child)
3. Run `PreRun` hook (this command only)
4. Run `Action` (this command only)
5. Run `PostRun` hook (this command only)
6. Run `PersistentPostRun` hooks (child → parent)

See [Lifecycle Hooks](lifecycle.md) for details.

## Best Practices

### Command Naming

- Use lowercase names
- Use hyphens for multi-word commands: `get-status`
- Keep names short and descriptive
- Use verbs for actions: `deploy`, `start`, `stop`

### Command Organization

Group related commands:

```go
// Good: Organized hierarchy
myapp database migrate
myapp database rollback
myapp database backup

// Bad: Flat structure
myapp db-migrate
myapp db-rollback
myapp db-backup
```

### Error Handling

Return errors from actions:

```go
cmd.Action(func(ctx context.Context, cmd *cli.Command) error {
    if err := validateConfig(); err != nil {
        return fmt.Errorf("invalid config: %w", err)
    }
    
    if err := deploy(); err != nil {
        return fmt.Errorf("deployment failed: %w", err)
    }
    
    return nil
})
```

## Examples

See the [examples/](../examples/) directory for:
- [basic](../examples/basic/) - Simple commands
- [subcommands](../examples/subcommands/) - Command hierarchies
- [lifecycle](../examples/lifecycle/) - Hook execution

## Next Steps

- Learn about [Flags and Arguments](flags-and-arguments.md)
- Understand [Lifecycle Hooks](lifecycle.md)
- Explore [API Reference](api-reference.md)
