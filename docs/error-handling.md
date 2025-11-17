# Error Handling

Best practices for error handling in CLI applications built with nyxstack/cli.

## Custom Error Types

nyxstack/cli provides three custom error types:

### CommandNotFoundError

Returned when a subcommand is not found:

```go
type CommandNotFoundError struct {
    Name string
}

func (e *CommandNotFoundError) Error() string {
    return fmt.Sprintf("command not found: %s", e.Name)
}
```

Usage:
```go
if err := app.Execute(); err != nil {
    if cmdErr, ok := err.(*cli.CommandNotFoundError); ok {
        fmt.Printf("Unknown command: %s\n", cmdErr.Name)
        fmt.Println("Run 'myapp --help' for usage")
        os.Exit(1)
    }
}
```

### ArgumentError

Returned when there are issues with arguments:

```go
type ArgumentError struct {
    Message string
}

func (e *ArgumentError) Error() string {
    return fmt.Sprintf("argument error: %s", e.Message)
}
```

Common cases:
- Missing required argument
- Too many arguments
- Invalid argument type

### FlagError

Returned when there are issues with flags:

```go
type FlagError struct {
    Message string
}

func (e *FlagError) Error() string {
    return fmt.Sprintf("flag error: %s", e.Message)
}
```

Common cases:
- Missing required flag
- Unknown flag
- Invalid flag value

## Action Error Handling

### Return Errors from Actions

```go
cmd.Action(func(ctx context.Context, cmd *cli.Command) error {
    // Simple error
    if !isValid() {
        return errors.New("validation failed")
    }
    
    // Wrapped error
    if err := doWork(); err != nil {
        return fmt.Errorf("work failed: %w", err)
    }
    
    return nil
})
```

### Error Wrapping

Use `%w` to wrap errors for better error chains:

```go
cmd.Action(func(ctx context.Context, cmd *cli.Command) error {
    if err := connectDB(); err != nil {
        return fmt.Errorf("database connection failed: %w", err)
    }
    
    if err := runMigrations(); err != nil {
        return fmt.Errorf("migrations failed: %w", err)
    }
    
    return nil
})
```

Check wrapped errors:
```go
if err := app.Execute(); err != nil {
    if errors.Is(err, sql.ErrNoRows) {
        fmt.Println("No data found")
    }
}
```

## Validation Errors

Use PreRun hooks for validation:

```go
cmd.PreRun(func(ctx context.Context, cmd *cli.Command) error {
    // Validate flags
    if port < 1 || port > 65535 {
        return fmt.Errorf("invalid port: %d (must be 1-65535)", port)
    }
    
    // Validate files
    if _, err := os.Stat(configFile); err != nil {
        return fmt.Errorf("config file not found: %s", configFile)
    }
    
    // Validate combinations
    if useSSL && certFile == "" {
        return errors.New("--cert-file required when --ssl is enabled")
    }
    
    return nil
})
```

## Exit Codes

### Standard Exit Codes

```go
func main() {
    app := cli.Root("myapp").Action(run)
    
    if err := app.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        
        // Choose exit code based on error type
        switch err.(type) {
        case *cli.CommandNotFoundError:
            os.Exit(127)  // Command not found
        case *cli.ArgumentError:
            os.Exit(2)    // Invalid usage
        case *cli.FlagError:
            os.Exit(2)    // Invalid usage
        default:
            os.Exit(1)    // General error
        }
    }
    
    os.Exit(0)  // Success
}
```

### Custom Exit Codes

```go
type ExitCodeError struct {
    Code    int
    Message string
}

func (e *ExitCodeError) Error() string {
    return e.Message
}

cmd.Action(func(ctx context.Context, cmd *cli.Command) error {
    if err := deploy(); err != nil {
        return &ExitCodeError{
            Code:    10,
            Message: "deployment failed",
        }
    }
    return nil
})

// In main
if err := app.Execute(); err != nil {
    if exitErr, ok := err.(*ExitCodeError); ok {
        fmt.Fprintf(os.Stderr, "Error: %v\n", exitErr.Message)
        os.Exit(exitErr.Code)
    }
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(1)
}
```

## Error Messages

### User-Friendly Messages

```go
cmd.Action(func(ctx context.Context, cmd *cli.Command) error {
    if err := validateAPIKey(apiKey); err != nil {
        // Bad: technical error
        return fmt.Errorf("invalid key format: %w", err)
        
        // Good: user-friendly message
        return errors.New("Invalid API key. Get your key from https://example.com/dashboard")
    }
    return nil
})
```

### Provide Context

```go
if err := loadConfig(configFile); err != nil {
    // Bad: no context
    return err
    
    // Good: includes file path
    return fmt.Errorf("failed to load config from %s: %w", configFile, err)
}
```

### Suggest Solutions

```go
if _, err := exec.LookPath("docker"); err != nil {
    return errors.New("docker not found in PATH. Install from https://docker.com")
}
```

## Error Recovery

### Graceful Degradation

```go
cmd.Action(func(ctx context.Context, cmd *cli.Command) error {
    // Try to load cache
    cache, err := loadCache()
    if err != nil {
        // Don't fail, just log
        fmt.Fprintf(os.Stderr, "Warning: cache load failed: %v\n", err)
        cache = newCache()
    }
    
    // Continue with degraded functionality
    return doWork(cache)
})
```

### Retry Logic

```go
func retryableAction(ctx context.Context, cmd *cli.Command) error {
    maxRetries := 3
    
    for i := 0; i < maxRetries; i++ {
        err := doWork()
        if err == nil {
            return nil
        }
        
        // Check if retryable
        if !isRetryable(err) {
            return err
        }
        
        fmt.Fprintf(os.Stderr, "Attempt %d failed: %v. Retrying...\n", i+1, err)
        time.Sleep(time.Second * time.Duration(i+1))
    }
    
    return errors.New("max retries exceeded")
}
```

## Hook Error Propagation

Errors in hooks stop execution:

```go
cmd.
    PersistentPreRun(func(ctx context.Context, cmd *cli.Command) error {
        if err := connectDB(); err != nil {
            return err  // Stops here, Action never runs
        }
        return nil
    }).
    Action(func(ctx context.Context, cmd *cli.Command) error {
        // Only runs if PersistentPreRun succeeded
        return nil
    })
```

### Error Cleanup

Even if Action fails, PersistentPostRun still runs:

```go
cmd.
    PersistentPreRun(func(ctx context.Context, cmd *cli.Command) error {
        db, _ := connectDB()
        ctx = context.WithValue(ctx, "db", db)
        return nil
    }).
    Action(func(ctx context.Context, cmd *cli.Command) error {
        return errors.New("action failed")  // Action fails
    }).
    PersistentPostRun(func(ctx context.Context, cmd *cli.Command) error {
        // Still runs! Clean up resources
        if db := ctx.Value("db"); db != nil {
            db.(*Database).Close()
        }
        return nil
    })
```

## Logging Errors

### Structured Logging

```go
import "log/slog"

cmd.Action(func(ctx context.Context, cmd *cli.Command) error {
    if err := doWork(); err != nil {
        slog.Error("work failed",
            "error", err,
            "command", cmd.GetName(),
            "user", os.Getenv("USER"),
        )
        return err
    }
    return nil
})
```

### Error Tracking

```go
import "github.com/getsentry/sentry-go"

cmd.Action(func(ctx context.Context, cmd *cli.Command) error {
    if err := doWork(); err != nil {
        sentry.CaptureException(err)
        return err
    }
    return nil
})
```

## Testing Error Handling

### Expected Errors

```go
func TestExpectedError(t *testing.T) {
    cmd := cli.Cmd("test").
        Action(func(ctx context.Context, cmd *cli.Command) error {
            return errors.New("expected error")
        })
    
    err := cmd.Execute()
    if err == nil {
        t.Error("expected error, got nil")
    }
    
    if err.Error() != "expected error" {
        t.Errorf("wrong error message: %v", err)
    }
}
```

### Error Types

```go
func TestErrorType(t *testing.T) {
    cmd := cli.Cmd("test").
        Arg("required", "Required arg", true)
    
    err := cmd.ExecuteWithArgs([]string{})
    
    if _, ok := err.(*cli.ArgumentError); !ok {
        t.Errorf("expected ArgumentError, got %T", err)
    }
}
```

### Error Wrapping

```go
func TestErrorWrapping(t *testing.T) {
    originalErr := errors.New("original")
    
    cmd := cli.Root("test").
        Action(func(ctx context.Context, cmd *cli.Command) error {
            return fmt.Errorf("wrapped: %w", originalErr)
        })
    
    err := cmd.Execute()
    
    if !errors.Is(err, originalErr) {
        t.Error("error chain broken")
    }
}
```

## Best Practices

### 1. Always Return Errors

```go
// Good
func (c *Command) Action(fn ActionFunc) error {
    return fn()
}

// Bad
func (c *Command) Action(fn ActionFunc) {
    fn()  // Swallows errors
}
```

### 2. Use Error Wrapping

```go
// Good
return fmt.Errorf("database query failed: %w", err)

// Bad
return errors.New("database query failed")  // Loses original error
```

### 3. Provide Context

```go
// Good
return fmt.Errorf("failed to parse config at %s: %w", path, err)

// Bad
return err  // No context about what failed
```

### 4. Check Error Types

```go
// Good
if _, ok := err.(*cli.ArgumentError); ok {
    // Handle argument errors specially
}

// Bad
if strings.Contains(err.Error(), "argument") {
    // Brittle string matching
}
```

### 5. Clean Up on Error

```go
// Good
func doWork() error {
    f, err := os.Open("file.txt")
    if err != nil {
        return err
    }
    defer f.Close()  // Cleans up even on error
    
    return process(f)
}
```

## Examples

See [examples/error-handling/](../examples/error-handling/) for complete working examples.

## Next Steps

- Learn about [Lifecycle Hooks](lifecycle.md) for validation
- Read [Context Support](context.md) for cancellation errors
- Check [Testing](testing.md) for error testing patterns
