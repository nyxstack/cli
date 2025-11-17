package cli

import (
	"context"
	"errors"
	"testing"
	"time"
)

// TestIntegrationFullWorkflow tests a complete command workflow
func TestIntegrationFullWorkflow(t *testing.T) {
	var (
		verbose     bool
		timeout     int
		executed    string
		receivedEnv string
	)

	root := Root("myapp").
		Description("My application").
		Flag(&verbose, "verbose", "v", false, "Verbose output").
		PersistentPreRun(func(ctx context.Context, cmd *Command) error {
			if verbose {
				executed += "verbose-pre;"
			}
			return nil
		})

	deploy := Cmd("deploy").
		Description("Deploy the application").
		Flag(&timeout, "timeout", "t", 30, "Timeout in seconds").
		Arg("environment", "Target environment", true).
		PreRun(func(ctx context.Context, cmd *Command) error {
			executed += "deploy-pre;"
			return nil
		}).
		Action(func(ctx context.Context, cmd *Command, env string) error {
			executed += "deploy-action;"
			receivedEnv = env
			return nil
		}).
		PostRun(func(ctx context.Context, cmd *Command) error {
			executed += "deploy-post;"
			return nil
		})

	root.AddCommand(deploy)

	err := root.ExecuteWithArgs([]string{"--verbose", "deploy", "--timeout=60", "production"})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}

	expectedOrder := "verbose-pre;deploy-pre;deploy-action;deploy-post;"
	if executed != expectedOrder {
		t.Errorf("expected execution order %q, got %q", expectedOrder, executed)
	}

	if receivedEnv != "production" {
		t.Errorf("expected environment 'production', got %q", receivedEnv)
	}

	timeoutFlag := deploy.flags.GetFlag("timeout")
	if val, _ := timeoutFlag.GetValue().(int); val != 60 {
		t.Errorf("expected timeout 60, got %d", val)
	}
}

// TestIntegrationErrorPropagation tests error handling through the stack
func TestIntegrationErrorPropagation(t *testing.T) {
	preRunErr := errors.New("prerun validation failed")

	root := Root("myapp")
	deploy := Cmd("deploy").
		PreRun(func(ctx context.Context, cmd *Command) error {
			return preRunErr
		}).
		Action(func(ctx context.Context, cmd *Command) error {
			t.Error("action should not run after PreRun error")
			return nil
		})

	root.AddCommand(deploy)

	err := root.ExecuteWithArgs([]string{"deploy"})
	if err == nil {
		t.Fatal("expected error from PreRun")
	}
	if err != preRunErr {
		t.Errorf("expected preRunErr, got %v", err)
	}
}

// TestIntegrationContextPropagation tests context values through execution
func TestIntegrationContextPropagation(t *testing.T) {
	type contextKey string
	key := contextKey("requestID")

	var receivedID string

	root := Root("myapp").
		Action(func(ctx context.Context, cmd *Command) error {
			if val := ctx.Value(key); val != nil {
				receivedID = val.(string)
			}
			return nil
		})

	ctx := context.WithValue(context.Background(), key, "req-123")
	err := root.execute(ctx, []string{})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}

	if receivedID != "req-123" {
		t.Errorf("expected requestID 'req-123', got %q", receivedID)
	}
}

// TestIntegrationContextCancellation tests context cancellation
func TestIntegrationContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	root := Root("myapp").
		Action(func(ctx context.Context, cmd *Command) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return nil
			}
		})

	err := root.execute(ctx, []string{})
	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

// TestIntegrationContextTimeout tests context timeout
func TestIntegrationContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	root := Root("myapp").
		Action(func(ctx context.Context, cmd *Command) error {
			time.Sleep(10 * time.Millisecond)
			return ctx.Err()
		})

	err := root.execute(ctx, []string{})
	if err != context.DeadlineExceeded {
		t.Errorf("expected context.DeadlineExceeded, got %v", err)
	}
}

// TestIntegrationComplexFlagParsing tests complex flag scenarios
func TestIntegrationComplexFlagParsing(t *testing.T) {
	var (
		verbose bool
		count   int
		tags    []string
		timeout time.Duration
		force   bool
	)

	root := Root("myapp").
		Flag(&verbose, "verbose", "v", false, "Verbose").
		Flag(&count, "count", "c", 5, "Count")

	deploy := Cmd("deploy").
		Flag(&tags, "tag", "t", []string{}, "Tags").
		Flag(&timeout, "timeout", "", 30*time.Second, "Timeout").
		Flag(&force, "force", "f", false, "Force").
		Action(func(ctx context.Context, cmd *Command) error {
			return nil
		})

	root.AddCommand(deploy)

	err := root.ExecuteWithArgs([]string{
		"-v",
		"--count=10",
		"deploy",
		"--tag=v1.0",
		"-t=v2.0",
		"--timeout=1m",
		"-f",
	})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}

	// Verify root flags
	if !verbose {
		t.Error("verbose should be true")
	}
	if count != 10 {
		t.Errorf("count should be 10, got %d", count)
	}

	// Verify deploy flags
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if tags[0] != "v1.0" || tags[1] != "v2.0" {
		t.Errorf("unexpected tags: %v", tags)
	}
	if timeout != 1*time.Minute {
		t.Errorf("timeout should be 1m, got %v", timeout)
	}
	if !force {
		t.Error("force should be true")
	}
}

// TestIntegrationNestedCommands tests deeply nested command structure
func TestIntegrationNestedCommands(t *testing.T) {
	executionPath := ""

	root := Root("kubectl").
		PersistentPreRun(func(ctx context.Context, cmd *Command) error {
			executionPath += "root-"
			return nil
		})

	get := Cmd("get").
		PersistentPreRun(func(ctx context.Context, cmd *Command) error {
			executionPath += "get-"
			return nil
		})

	pods := Cmd("pods").
		Action(func(ctx context.Context, cmd *Command) error {
			executionPath += "pods"
			return nil
		})

	root.AddCommand(get)
	get.AddCommand(pods)

	err := root.ExecuteWithArgs([]string{"get", "pods"})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}

	expected := "root-get-pods"
	if executionPath != expected {
		t.Errorf("expected path %q, got %q", expected, executionPath)
	}
}

// TestIntegrationMixedArgumentTypes tests various argument types together
func TestIntegrationMixedArgumentTypes(t *testing.T) {
	var (
		name     string
		count    int
		ratio    float64
		enabled  bool
		duration time.Duration
	)

	root := Root("myapp").
		Arg("name", "Name", true).
		Arg("count", "Count", true).
		Arg("ratio", "Ratio", true).
		Arg("enabled", "Enabled", true).
		Arg("duration", "Duration", true).
		Action(func(ctx context.Context, cmd *Command, n string, c int, r float64, e bool, d time.Duration) error {
			name = n
			count = c
			ratio = r
			enabled = e
			duration = d
			return nil
		})

	err := root.ExecuteWithArgs([]string{"test", "42", "3.14", "true", "5s"})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}

	if name != "test" {
		t.Errorf("name: expected 'test', got %q", name)
	}
	if count != 42 {
		t.Errorf("count: expected 42, got %d", count)
	}
	if ratio != 3.14 {
		t.Errorf("ratio: expected 3.14, got %f", ratio)
	}
	if !enabled {
		t.Error("enabled: expected true, got false")
	}
	if duration != 5*time.Second {
		t.Errorf("duration: expected 5s, got %v", duration)
	}
}

// TestIntegrationValidationWorkflow tests validation in PreRun hooks
func TestIntegrationValidationWorkflow(t *testing.T) {
	root := Root("myapp")

	deploy := Cmd("deploy").
		Arg("environment", "Environment", true).
		PreRun(func(ctx context.Context, cmd *Command) error {
			// This would normally validate the environment
			return &ArgumentError{
				Arg: "environment",
				Msg: "invalid environment",
				Cmd: cmd,
			}
		}).
		Action(func(ctx context.Context, cmd *Command, env string) error {
			t.Error("action should not run after validation error")
			return nil
		})

	root.AddCommand(deploy)

	err := root.ExecuteWithArgs([]string{"deploy", "invalid-env"})
	if err == nil {
		t.Fatal("expected validation error")
	}

	argErr, ok := err.(*ArgumentError)
	if !ok {
		t.Errorf("expected ArgumentError, got %T", err)
	}
	if argErr.Arg != "environment" {
		t.Errorf("expected error for 'environment', got %q", argErr.Arg)
	}
}

// TestIntegrationFlagRequired tests required flag validation
func TestIntegrationFlagRequired(t *testing.T) {
	var apiKey string

	root := Root("myapp").
		FlagRequired(&apiKey, "api-key", "", "", "API key").
		Action(func(ctx context.Context, cmd *Command) error {
			return nil
		})

	// Should fail without required flag
	err := root.ExecuteWithArgs([]string{})
	if err == nil {
		t.Fatal("expected error for missing required flag")
	}

	flagErr, ok := err.(*FlagError)
	if !ok {
		t.Errorf("expected FlagError, got %T", err)
	}
	if flagErr.Flag != "api-key" {
		t.Errorf("expected error for 'api-key', got %q", flagErr.Flag)
	}

	// Should succeed with required flag
	err = root.ExecuteWithArgs([]string{"--api-key=secret123"})
	if err != nil {
		t.Errorf("should succeed with required flag, got %v", err)
	}
	if apiKey != "secret123" {
		t.Errorf("expected apiKey 'secret123', got %q", apiKey)
	}
}
