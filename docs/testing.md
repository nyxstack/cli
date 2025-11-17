# Testing Your CLI Applications

Guide to testing CLI applications built with nyxstack/cli.

## Testing Commands

### Basic Test

```go
func TestMyCommand(t *testing.T) {
    var executed bool
    
    cmd := cli.Root("test").
        Action(func(ctx context.Context, cmd *cli.Command) error {
            executed = true
            return nil
        })
    
    err := cmd.Execute()
    
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    
    if !executed {
        t.Error("command was not executed")
    }
}
```

### Testing with Arguments

Use `ExecuteWithArgs()` to test with custom arguments:

```go
func TestDeployCommand(t *testing.T) {
    var env string
    var replicas int
    
    cmd := cli.Cmd("deploy").
        Arg("environment", "Target env", true).
        Arg("replicas", "Number of replicas", true).
        Action(func(ctx context.Context, cmd *cli.Command, e string, r int) error {
            env = e
            replicas = r
            return nil
        })
    
    err := cmd.ExecuteWithArgs([]string{"production", "3"})
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if env != "production" {
        t.Errorf("expected env=production, got %s", env)
    }
    
    if replicas != 3 {
        t.Errorf("expected replicas=3, got %d", replicas)
    }
}
```

### Testing with Flags

```go
func TestFlagParsing(t *testing.T) {
    var verbose bool
    var port int
    
    cmd := cli.Root("test").
        Flag(&verbose, "verbose", "v", false, "Verbose").
        Flag(&port, "port", "p", 8080, "Port").
        Action(func(ctx context.Context, cmd *cli.Command) error {
            return nil
        })
    
    err := cmd.ExecuteWithArgs([]string{"--verbose=true", "--port=3000"})
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if !verbose {
        t.Error("verbose should be true")
    }
    
    if port != 3000 {
        t.Errorf("expected port=3000, got %d", port)
    }
}
```

### Testing Subcommands

```go
func TestSubcommand(t *testing.T) {
    var executed string
    
    root := cli.Root("myapp")
    
    deploy := cli.Cmd("deploy").
        Action(func(ctx context.Context, cmd *cli.Command) error {
            executed = "deploy"
            return nil
        })
    
    rollback := cli.Cmd("rollback").
        Action(func(ctx context.Context, cmd *cli.Command) error {
            executed = "rollback"
            return nil
        })
    
    root.AddCommand(deploy)
    root.AddCommand(rollback)
    
    // Test deploy
    err := root.ExecuteWithArgs([]string{"deploy"})
    if err != nil || executed != "deploy" {
        t.Errorf("deploy command failed")
    }
    
    // Test rollback
    err = root.ExecuteWithArgs([]string{"rollback"})
    if err != nil || executed != "rollback" {
        t.Errorf("rollback command failed")
    }
}
```

## Testing Error Handling

### Expected Errors

```go
func TestMissingRequiredArg(t *testing.T) {
    cmd := cli.Cmd("deploy").
        Arg("environment", "Target env", true).
        Action(func(ctx context.Context, cmd *cli.Command, env string) error {
            return nil
        })
    
    err := cmd.ExecuteWithArgs([]string{})
    
    if err == nil {
        t.Error("expected error for missing argument")
    }
    
    // Check error type
    if _, ok := err.(*cli.ArgumentError); !ok {
        t.Errorf("expected ArgumentError, got %T", err)
    }
}
```

### Custom Error Validation

```go
func TestCustomError(t *testing.T) {
    expectedErr := errors.New("custom error")
    
    cmd := cli.Root("test").
        Action(func(ctx context.Context, cmd *cli.Command) error {
            return expectedErr
        })
    
    err := cmd.Execute()
    
    if !errors.Is(err, expectedErr) {
        t.Errorf("expected error %v, got %v", expectedErr, err)
    }
}
```

## Testing Context

### Context Cancellation

```go
func TestContextCancellation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    
    cmd := cli.Root("test").
        Action(func(ctx context.Context, cmd *cli.Command) error {
            <-ctx.Done()
            return ctx.Err()
        })
    
    // Cancel immediately
    cancel()
    
    err := cmd.ExecuteContext(ctx)
    
    if !errors.Is(err, context.Canceled) {
        t.Errorf("expected context.Canceled, got %v", err)
    }
}
```

### Context Timeout

```go
func TestContextTimeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
    defer cancel()
    
    cmd := cli.Root("test").
        Action(func(ctx context.Context, cmd *cli.Command) error {
            time.Sleep(100 * time.Millisecond)
            return nil
        })
    
    err := cmd.ExecuteContext(ctx)
    
    if !errors.Is(err, context.DeadlineExceeded) {
        t.Errorf("expected context.DeadlineExceeded, got %v", err)
    }
}
```

## Testing Lifecycle Hooks

```go
func TestLifecycleOrder(t *testing.T) {
    var order []string
    
    cmd := cli.Root("test").
        PersistentPreRun(func(ctx context.Context, cmd *cli.Command) error {
            order = append(order, "persistent-pre")
            return nil
        }).
        PreRun(func(ctx context.Context, cmd *cli.Command) error {
            order = append(order, "pre")
            return nil
        }).
        Action(func(ctx context.Context, cmd *cli.Command) error {
            order = append(order, "action")
            return nil
        }).
        PostRun(func(ctx context.Context, cmd *cli.Command) error {
            order = append(order, "post")
            return nil
        }).
        PersistentPostRun(func(ctx context.Context, cmd *cli.Command) error {
            order = append(order, "persistent-post")
            return nil
        })
    
    err := cmd.Execute()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    expected := []string{"persistent-pre", "pre", "action", "post", "persistent-post"}
    if !reflect.DeepEqual(order, expected) {
        t.Errorf("expected order %v, got %v", expected, order)
    }
}
```

## Testing Flag Inheritance

```go
func TestFlagInheritance(t *testing.T) {
    var verbose bool
    
    root := cli.Root("test").
        Flag(&verbose, "verbose", "v", false, "Verbose")
    
    child := cli.Cmd("child").
        Action(func(ctx context.Context, cmd *cli.Command) error {
            if !verbose {
                t.Error("child should have access to parent's verbose flag")
            }
            return nil
        })
    
    root.AddCommand(child)
    
    err := root.ExecuteWithArgs([]string{"--verbose=true", "child"})
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}
```

## Table-Driven Tests

```go
func TestMultipleScenarios(t *testing.T) {
    tests := []struct {
        name    string
        args    []string
        wantErr bool
    }{
        {"valid args", []string{"production", "3"}, false},
        {"missing arg", []string{"production"}, true},
        {"invalid type", []string{"production", "invalid"}, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cmd := cli.Cmd("deploy").
                Arg("environment", "Env", true).
                Arg("replicas", "Replicas", true).
                Action(func(ctx context.Context, cmd *cli.Command, env string, replicas int) error {
                    return nil
                })
            
            err := cmd.ExecuteWithArgs(tt.args)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("wantErr=%v, got error=%v", tt.wantErr, err)
            }
        })
    }
}
```

## Mocking and Test Helpers

### Test Helper Function

```go
func executeCommand(t *testing.T, cmd *cli.Command, args []string) error {
    t.Helper()
    
    err := cmd.ExecuteWithArgs(args)
    return err
}

func TestWithHelper(t *testing.T) {
    cmd := cli.Root("test")
    err := executeCommand(t, cmd, []string{"--help"})
    // ...
}
```

### Command Builder for Tests

```go
func testCommand(action func(context.Context, *cli.Command) error) *cli.Command {
    return cli.Root("test").Action(action)
}

func TestBuilder(t *testing.T) {
    cmd := testCommand(func(ctx context.Context, cmd *cli.Command) error {
        return nil
    })
    
    err := cmd.Execute()
    if err != nil {
        t.Error("command failed")
    }
}
```

## Best Practices

1. **Use `ExecuteWithArgs()` for tests**
   ```go
   err := cmd.ExecuteWithArgs([]string{"--flag=value", "arg1"})
   ```

2. **Test both success and failure cases**
   ```go
   t.Run("success", func(t *testing.T) { /* ... */ })
   t.Run("failure", func(t *testing.T) { /* ... */ })
   ```

3. **Use table-driven tests for multiple scenarios**
   ```go
   tests := []struct{ name, args, want }{ /* ... */ }
   ```

4. **Test error types, not just error messages**
   ```go
   if _, ok := err.(*cli.ArgumentError); !ok { /* ... */ }
   ```

5. **Test context propagation**
   ```go
   ctx := context.WithValue(context.Background(), "key", "value")
   cmd.ExecuteContext(ctx)
   ```

## Running Tests

```bash
# Run all tests
go test -v

# Run specific test
go test -v -run TestMyCommand

# Run with coverage
go test -v -cover

# Run with race detector
go test -v -race
```

## Next Steps

- See [TESTING.md](../TESTING.md) for framework test coverage
- Check [Examples](../examples/) for more patterns
- Read [API Reference](api-reference.md) for details
