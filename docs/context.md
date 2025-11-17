# Context Support

Using context.Context for cancellation, timeouts, and value propagation.

## Basic Usage

Every action receives a context as the first parameter:

```go
cmd := cli.Root("myapp").
    Action(func(ctx context.Context, cmd *cli.Command) error {
        // ctx is available for all operations
        return doWork(ctx)
    })
```

## Cancellation

### Basic Cancellation

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

app := cli.Root("myapp").
    Action(func(ctx context.Context, cmd *cli.Command) error {
        // Check for cancellation
        select {
        case <-ctx.Done():
            return ctx.Err()  // context.Canceled
        default:
            return doWork(ctx)
        }
    })

err := app.ExecuteContext(ctx)
```

### Graceful Shutdown

```go
func main() {
    // Handle interrupt signal
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        <-sigChan
        fmt.Println("\nShutting down gracefully...")
        cancel()
    }()
    
    app := cli.Root("myapp").
        Action(func(ctx context.Context, cmd *cli.Command) error {
            // Long-running operation
            for {
                select {
                case <-ctx.Done():
                    fmt.Println("Operation cancelled")
                    return ctx.Err()
                case <-time.After(1 * time.Second):
                    fmt.Print(".")
                }
            }
        })
    
    if err := app.ExecuteContext(ctx); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    }
}
```

## Timeouts

### Command-Level Timeout

```go
func main() {
    // 30 second timeout for entire command
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    app := cli.Root("myapp").
        Action(func(ctx context.Context, cmd *cli.Command) error {
            return doWorkWithTimeout(ctx)
        })
    
    if err := app.ExecuteContext(ctx); err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            fmt.Println("Command timed out!")
        }
    }
}
```

### Operation-Level Timeout

```go
cmd.Action(func(ctx context.Context, cmd *cli.Command) error {
    // Set timeout for specific operation
    opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    result, err := callAPI(opCtx)
    if errors.Is(err, context.DeadlineExceeded) {
        return fmt.Errorf("API call timed out")
    }
    
    return processResult(result)
})
```

## Value Propagation

### Passing Values Through Context

```go
root := cli.Root("myapp").
    PersistentPreRun(func(ctx context.Context, cmd *cli.Command) error {
        // Load configuration
        config := loadConfig()
        
        // Store in context
        ctx = context.WithValue(ctx, "config", config)
        
        // Initialize database
        db, err := connectDB(config.DBUrl)
        if err != nil {
            return err
        }
        
        ctx = context.WithValue(ctx, "db", db)
        return nil
    }).
    Action(func(ctx context.Context, cmd *cli.Command) error {
        // Retrieve from context
        config := ctx.Value("config").(*Config)
        db := ctx.Value("db").(*Database)
        
        return doWork(config, db)
    }).
    PersistentPostRun(func(ctx context.Context, cmd *cli.Command) error {
        // Clean up
        if db := ctx.Value("db"); db != nil {
            db.(*Database).Close()
        }
        return nil
    })
```

### Type-Safe Context Keys

```go
// Define context keys
type contextKey string

const (
    configKey contextKey = "config"
    dbKey     contextKey = "db"
)

// Helper functions
func withConfig(ctx context.Context, cfg *Config) context.Context {
    return context.WithValue(ctx, configKey, cfg)
}

func getConfig(ctx context.Context) *Config {
    if cfg, ok := ctx.Value(configKey).(*Config); ok {
        return cfg
    }
    return nil
}

// Usage
root.PersistentPreRun(func(ctx context.Context, cmd *cli.Command) error {
    config := loadConfig()
    ctx = withConfig(ctx, config)
    return nil
})

cmd.Action(func(ctx context.Context, cmd *cli.Command) error {
    config := getConfig(ctx)
    // Use config...
    return nil
})
```

## HTTP Client with Context

```go
cmd.Action(func(ctx context.Context, cmd *cli.Command) error {
    client := &http.Client{
        Timeout: 10 * time.Second,
    }
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return err
    }
    
    resp, err := client.Do(req)
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            return fmt.Errorf("request timed out")
        }
        if errors.Is(err, context.Canceled) {
            return fmt.Errorf("request cancelled")
        }
        return err
    }
    defer resp.Body.Close()
    
    // Process response...
    return nil
})
```

## Database Operations

```go
cmd.Action(func(ctx context.Context, cmd *cli.Command) error {
    db := getDB(ctx)
    
    // Query with context (auto-cancels if ctx cancelled)
    rows, err := db.QueryContext(ctx, "SELECT * FROM users")
    if err != nil {
        if errors.Is(err, context.Canceled) {
            return fmt.Errorf("query cancelled")
        }
        return err
    }
    defer rows.Close()
    
    // Process rows...
    return nil
})
```

## Goroutines with Context

### Background Worker

```go
cmd.Action(func(ctx context.Context, cmd *cli.Command) error {
    // Start background worker
    go func() {
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()
        
        for {
            select {
            case <-ctx.Done():
                fmt.Println("Worker stopped")
                return
            case <-ticker.C:
                fmt.Println("Working...")
            }
        }
    }()
    
    // Wait for signal or completion
    <-ctx.Done()
    return ctx.Err()
})
```

### Multiple Goroutines

```go
cmd.Action(func(ctx context.Context, cmd *cli.Command) error {
    var wg sync.WaitGroup
    errChan := make(chan error, 3)
    
    // Launch 3 workers
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            select {
            case <-ctx.Done():
                errChan <- fmt.Errorf("worker %d cancelled", id)
                return
            default:
                if err := doWork(ctx, id); err != nil {
                    errChan <- err
                }
            }
        }(i)
    }
    
    // Wait for completion
    wg.Wait()
    close(errChan)
    
    // Check for errors
    for err := range errChan {
        if err != nil {
            return err
        }
    }
    
    return nil
})
```

## Best Practices

### 1. Always Pass Context

```go
// Good
func doWork(ctx context.Context) error {
    // Can check ctx.Done()
}

// Bad
func doWork() error {
    // No way to cancel
}
```

### 2. Don't Store Context in Structs

```go
// Bad
type Server struct {
    ctx context.Context  // Don't do this
}

// Good
func (s *Server) Start(ctx context.Context) error {
    // Pass as parameter
}
```

### 3. Check for Cancellation Early

```go
func doWork(ctx context.Context) error {
    // Check before expensive operation
    if ctx.Err() != nil {
        return ctx.Err()
    }
    
    // Do work...
    return nil
}
```

### 4. Use Type-Safe Keys

```go
// Good
type contextKey string
const userKey contextKey = "user"

// Bad
ctx.Value("user")  // String keys can collide
```

### 5. Set Reasonable Timeouts

```go
// Good
ctx, cancel := context.WithTimeout(ctx, 30*time.Second)

// Bad
ctx, cancel := context.WithTimeout(ctx, 1*time.Hour)  // Too long
```

## Testing with Context

```go
func TestWithContext(t *testing.T) {
    ctx := context.Background()
    
    cmd := cli.Root("test").
        Action(func(ctx context.Context, cmd *cli.Command) error {
            // Test action with context
            return nil
        })
    
    err := cmd.ExecuteContext(ctx)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
}

func TestCancellation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    cancel()  // Cancel immediately
    
    cmd := cli.Root("test").
        Action(func(ctx context.Context, cmd *cli.Command) error {
            select {
            case <-ctx.Done():
                return ctx.Err()
            default:
                return errors.New("should have been cancelled")
            }
        })
    
    err := cmd.ExecuteContext(ctx)
    if !errors.Is(err, context.Canceled) {
        t.Error("expected context.Canceled error")
    }
}
```

## Common Patterns

### Retry with Context

```go
func retryWithContext(ctx context.Context, maxRetries int, fn func(context.Context) error) error {
    for i := 0; i < maxRetries; i++ {
        // Check cancellation
        if ctx.Err() != nil {
            return ctx.Err()
        }
        
        err := fn(ctx)
        if err == nil {
            return nil
        }
        
        // Wait before retry (with context awareness)
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(time.Second * time.Duration(i+1)):
            continue
        }
    }
    return fmt.Errorf("max retries exceeded")
}
```

### Fan-Out/Fan-In

```go
func fanOutFanIn(ctx context.Context, items []Item) ([]Result, error) {
    resultChan := make(chan Result, len(items))
    errChan := make(chan error, 1)
    
    // Fan out
    for _, item := range items {
        go func(i Item) {
            result, err := processItem(ctx, i)
            if err != nil {
                select {
                case errChan <- err:
                default:
                }
                return
            }
            resultChan <- result
        }(item)
    }
    
    // Fan in
    var results []Result
    for i := 0; i < len(items); i++ {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        case err := <-errChan:
            return nil, err
        case result := <-resultChan:
            results = append(results, result)
        }
    }
    
    return results, nil
}
```

## Examples

See [examples/context-usage/](../examples/context-usage/) for complete working examples.

## Next Steps

- Learn about [Lifecycle Hooks](lifecycle.md)
- Read [Error Handling](error-handling.md)
- Check [Testing](testing.md) for context tests
