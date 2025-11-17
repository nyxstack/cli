# Lifecycle Hooks

Control the execution flow of your commands with lifecycle hooks.

## Hook Types

nyxstack/cli provides four lifecycle hooks:

1. **PersistentPreRun** - Runs before action, for this command and all children
2. **PreRun** - Runs before action, only for this command
3. **Action** - Main command logic
4. **PostRun** - Runs after action, only for this command
5. **PersistentPostRun** - Runs after action, for this command and all children

## Execution Order

```
PersistentPreRun (parent)
    ↓
PersistentPreRun (child)
    ↓
PreRun (child only)
    ↓
Action (child only)
    ↓
PostRun (child only)
    ↓
PersistentPostRun (child)
    ↓
PersistentPostRun (parent)
```

## Basic Usage

```go
cmd := cli.Root("myapp").
    PersistentPreRun(func(ctx context.Context, cmd *cli.Command) error {
        fmt.Println("1. PersistentPreRun")
        return nil
    }).
    PreRun(func(ctx context.Context, cmd *cli.Command) error {
        fmt.Println("2. PreRun")
        return nil
    }).
    Action(func(ctx context.Context, cmd *cli.Command) error {
        fmt.Println("3. Action")
        return nil
    }).
    PostRun(func(ctx context.Context, cmd *cli.Command) error {
        fmt.Println("4. PostRun")
        return nil
    }).
    PersistentPostRun(func(ctx context.Context, cmd *cli.Command) error {
        fmt.Println("5. PersistentPostRun")
        return nil
    })
```

Output:
```
1. PersistentPreRun
2. PreRun
3. Action
4. PostRun
5. PersistentPostRun
```

## PersistentPreRun

Use for initialization that applies to this command and all subcommands:

```go
root := cli.Root("myapp").
    PersistentPreRun(func(ctx context.Context, cmd *cli.Command) error {
        // Initialize database connection
        db, err := connectDatabase()
        if err != nil {
            return fmt.Errorf("database connection failed: %w", err)
        }
        
        // Store in context for child commands
        ctx = context.WithValue(ctx, "db", db)
        
        // Initialize logging
        setupLogger()
        
        return nil
    })
```

This runs before **every** subcommand under `root`.

## PreRun

Use for validation or setup specific to one command:

```go
deploy := cli.Cmd("deploy").
    PreRun(func(ctx context.Context, cmd *cli.Command) error {
        // Validate environment
        if !isValidEnvironment(env) {
            return fmt.Errorf("invalid environment: %s", env)
        }
        
        // Check prerequisites
        if !hasRequiredTools() {
            return errors.New("missing required tools")
        }
        
        return nil
    }).
    Action(func(ctx context.Context, cmd *cli.Command) error {
        // Deploy logic here
        return nil
    })
```

This runs only for the `deploy` command, not its children.

## Action

The main command logic:

```go
cmd := cli.Cmd("build").
    Action(func(ctx context.Context, cmd *cli.Command) error {
        // Build application
        fmt.Println("Building...")
        
        if err := compile(); err != nil {
            return fmt.Errorf("compilation failed: %w", err)
        }
        
        fmt.Println("Build successful!")
        return nil
    })
```

## PostRun

Use for cleanup or finalization:

```go
cmd := cli.Cmd("test").
    Action(func(ctx context.Context, cmd *cli.Command) error {
        return runTests()
    }).
    PostRun(func(ctx context.Context, cmd *cli.Command) error {
        // Generate test report
        generateReport()
        
        // Clean up temporary files
        cleanupTempFiles()
        
        return nil
    })
```

Runs only if Action succeeds.

## PersistentPostRun

Use for cleanup that applies to this command and all subcommands:

```go
root := cli.Root("myapp").
    PersistentPreRun(func(ctx context.Context, cmd *cli.Command) error {
        db, _ := connectDatabase()
        ctx = context.WithValue(ctx, "db", db)
        return nil
    }).
    PersistentPostRun(func(ctx context.Context, cmd *cli.Command) error {
        // Close database connection
        db := ctx.Value("db").(*Database)
        db.Close()
        
        // Flush logs
        flushLogs()
        
        return nil
    })
```

Runs after **every** subcommand completes.

## Error Handling

If any hook returns an error, execution stops immediately:

```go
cmd := cli.Root("myapp").
    PersistentPreRun(func(ctx context.Context, cmd *cli.Command) error {
        if !checkLicense() {
            return errors.New("invalid license")  // Stops here
        }
        return nil
    }).
    Action(func(ctx context.Context, cmd *cli.Command) error {
        // Never runs if PreRun fails
        return nil
    })
```

## Inheritance Example

```go
root := cli.Root("myapp").
    PersistentPreRun(func(ctx context.Context, cmd *cli.Command) error {
        fmt.Println("Root: PersistentPreRun")
        return nil
    }).
    PersistentPostRun(func(ctx context.Context, cmd *cli.Command) error {
        fmt.Println("Root: PersistentPostRun")
        return nil
    })

database := cli.Cmd("database").
    PersistentPreRun(func(ctx context.Context, cmd *cli.Command) error {
        fmt.Println("Database: PersistentPreRun")
        return nil
    }).
    PersistentPostRun(func(ctx context.Context, cmd *cli.Command) error {
        fmt.Println("Database: PersistentPostRun")
        return nil
    })

migrate := cli.Cmd("migrate").
    PreRun(func(ctx context.Context, cmd *cli.Command) error {
        fmt.Println("Migrate: PreRun")
        return nil
    }).
    Action(func(ctx context.Context, cmd *cli.Command) error {
        fmt.Println("Migrate: Action")
        return nil
    }).
    PostRun(func(ctx context.Context, cmd *cli.Command) error {
        fmt.Println("Migrate: PostRun")
        return nil
    })

database.AddCommand(migrate)
root.AddCommand(database)
```

Running `myapp database migrate` produces:

```
Root: PersistentPreRun
Database: PersistentPreRun
Migrate: PreRun
Migrate: Action
Migrate: PostRun
Database: PersistentPostRun
Root: PersistentPostRun
```

## Common Patterns

### Database Connection

```go
root.PersistentPreRun(func(ctx context.Context, cmd *cli.Command) error {
    db, err := sql.Open("postgres", connString)
    if err != nil {
        return err
    }
    ctx = context.WithValue(ctx, "db", db)
    return nil
}).PersistentPostRun(func(ctx context.Context, cmd *cli.Command) error {
    if db := ctx.Value("db"); db != nil {
        db.(*sql.DB).Close()
    }
    return nil
})
```

### Authentication

```go
cmd.PreRun(func(ctx context.Context, cmd *cli.Command) error {
    token := os.Getenv("API_TOKEN")
    if token == "" {
        return errors.New("API_TOKEN not set")
    }
    
    if !validateToken(token) {
        return errors.New("invalid token")
    }
    
    return nil
})
```

### Configuration Loading

```go
root.PersistentPreRun(func(ctx context.Context, cmd *cli.Command) error {
    config, err := loadConfig(configFile)
    if err != nil {
        return fmt.Errorf("failed to load config: %w", err)
    }
    
    ctx = context.WithValue(ctx, "config", config)
    return nil
})
```

### Logging Setup

```go
root.PersistentPreRun(func(ctx context.Context, cmd *cli.Command) error {
    logLevel := "info"
    if verbose {
        logLevel = "debug"
    }
    
    logger := setupLogger(logLevel)
    ctx = context.WithValue(ctx, "logger", logger)
    return nil
}).PersistentPostRun(func(ctx context.Context, cmd *cli.Command) error {
    if logger := ctx.Value("logger"); logger != nil {
        logger.(*Logger).Sync()
    }
    return nil
})
```

### Validation

```go
cmd.PreRun(func(ctx context.Context, cmd *cli.Command) error {
    // Validate flags
    if port < 1 || port > 65535 {
        return fmt.Errorf("invalid port: %d", port)
    }
    
    // Check file exists
    if _, err := os.Stat(configFile); err != nil {
        return fmt.Errorf("config file not found: %s", configFile)
    }
    
    // Validate combination of flags
    if useSSL && certFile == "" {
        return errors.New("--cert-file required when using SSL")
    }
    
    return nil
})
```

## Best Practices

1. **Use PersistentPreRun for shared initialization**
   - Database connections
   - Configuration loading
   - Logging setup

2. **Use PreRun for validation**
   - Flag validation
   - File existence checks
   - Permission checks

3. **Keep Action focused**
   - Only core command logic
   - Assume PreRun validated everything

4. **Use PostRun for cleanup**
   - Generate reports
   - Clean temporary files
   - Update caches

5. **Use PersistentPostRun for shared cleanup**
   - Close database connections
   - Flush logs
   - Release resources

6. **Return meaningful errors**
   ```go
   // Good
   return fmt.Errorf("failed to connect to %s: %w", host, err)
   
   // Bad
   return err
   ```

## Examples

See the [examples/](../examples/) directory:
- [lifecycle](../examples/lifecycle/) - Complete lifecycle example
- [validation](../examples/validation/) - PreRun validation
- [error-handling](../examples/error-handling/) - Error propagation

## Next Steps

- Learn about [Context Support](context.md)
- Explore [Error Handling](error-handling.md)
- Read the [API Reference](api-reference.md)
