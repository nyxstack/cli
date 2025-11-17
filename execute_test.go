package cli

import (
	"context"
	"errors"
	"testing"
)

// TestExecuteBasic tests basic command execution
func TestExecuteBasic(t *testing.T) {
	executed := false
	cmd := Root("test").
		Action(func(ctx context.Context, c *Command) error {
			executed = true
			return nil
		})

	err := cmd.ExecuteWithArgs([]string{})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}
	if !executed {
		t.Error("action was not executed")
	}
}

// TestExecuteWithArgs tests ExecuteWithArgs
func TestExecuteWithArgs(t *testing.T) {
	var receivedArg string
	cmd := Root("test").
		Arg("name", "Name arg", true).
		Action(func(ctx context.Context, c *Command, name string) error {
			receivedArg = name
			return nil
		})

	err := cmd.ExecuteWithArgs([]string{"testname"})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}
	if receivedArg != "testname" {
		t.Errorf("expected 'testname', got %q", receivedArg)
	}
}

// TestExecuteContext tests ExecuteContext
func TestExecuteContext(t *testing.T) {
	type contextKey string
	key := contextKey("testkey")

	var receivedValue string
	cmd := Root("test").
		Action(func(ctx context.Context, c *Command) error {
			if val := ctx.Value(key); val != nil {
				receivedValue = val.(string)
			}
			return nil
		})

	ctx := context.WithValue(context.Background(), key, "testvalue")
	err := cmd.execute(ctx, []string{})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}
	if receivedValue != "testvalue" {
		t.Errorf("expected 'testvalue', got %q", receivedValue)
	}
}

// TestExecuteSubcommand tests subcommand execution
func TestExecuteSubcommand(t *testing.T) {
	rootExecuted := false
	childExecuted := false

	root := Root("app").
		Action(func(ctx context.Context, c *Command) error {
			rootExecuted = true
			return nil
		})

	child := Cmd("deploy").
		Action(func(ctx context.Context, c *Command) error {
			childExecuted = true
			return nil
		})

	root.AddCommand(child)

	err := root.ExecuteWithArgs([]string{"deploy"})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}
	if rootExecuted {
		t.Error("root action should not execute when subcommand is called")
	}
	if !childExecuted {
		t.Error("child action should execute")
	}
}

// TestExecuteArgumentTypes tests different argument types
func TestExecuteArgumentTypes(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		checkFunc func(*testing.T)
	}{
		{
			name: "string argument",
			args: []string{"hello"},
			checkFunc: func(t *testing.T) {
				var received string
				cmd := Root("test").
					Arg("msg", "Message", true).
					Action(func(ctx context.Context, c *Command, msg string) error {
						received = msg
						return nil
					})
				cmd.ExecuteWithArgs([]string{"hello"})
				if received != "hello" {
					t.Errorf("expected 'hello', got %q", received)
				}
			},
		},
		{
			name: "int argument",
			args: []string{"42"},
			checkFunc: func(t *testing.T) {
				var received int
				cmd := Root("test").
					Arg("count", "Count", true).
					Action(func(ctx context.Context, c *Command, count int) error {
						received = count
						return nil
					})
				cmd.ExecuteWithArgs([]string{"42"})
				if received != 42 {
					t.Errorf("expected 42, got %d", received)
				}
			},
		},
		{
			name: "bool argument",
			args: []string{"true"},
			checkFunc: func(t *testing.T) {
				var received bool
				cmd := Root("test").
					Arg("enabled", "Enabled", true).
					Action(func(ctx context.Context, c *Command, enabled bool) error {
						received = enabled
						return nil
					})
				cmd.ExecuteWithArgs([]string{"true"})
				if !received {
					t.Error("expected true, got false")
				}
			},
		},
		{
			name: "float64 argument",
			args: []string{"3.14"},
			checkFunc: func(t *testing.T) {
				var received float64
				cmd := Root("test").
					Arg("value", "Value", true).
					Action(func(ctx context.Context, c *Command, value float64) error {
						received = value
						return nil
					})
				cmd.ExecuteWithArgs([]string{"3.14"})
				if received != 3.14 {
					t.Errorf("expected 3.14, got %f", received)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.checkFunc(t)
		})
	}
}

// TestExecuteVariadicArgs tests variadic arguments
func TestExecuteVariadicArgs(t *testing.T) {
	var received []string
	cmd := Root("test").
		Action(func(ctx context.Context, c *Command, args ...string) error {
			received = args
			return nil
		})

	err := cmd.ExecuteWithArgs([]string{"one", "two", "three"})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}
	if len(received) != 3 {
		t.Fatalf("expected 3 args, got %d", len(received))
	}
	if received[0] != "one" || received[1] != "two" || received[2] != "three" {
		t.Errorf("unexpected args: %v", received)
	}
}

// TestExecuteVariadicArgsEmpty tests variadic with no arguments
func TestExecuteVariadicArgsEmpty(t *testing.T) {
	var received []string
	cmd := Root("test").
		Action(func(ctx context.Context, c *Command, args ...string) error {
			received = args
			return nil
		})

	err := cmd.ExecuteWithArgs([]string{})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}
	if len(received) != 0 {
		t.Errorf("expected 0 args, got %d", len(received))
	}
}

// TestExecuteErrorHandling tests error returns
func TestExecuteErrorHandling(t *testing.T) {
	expectedErr := errors.New("test error")
	cmd := Root("test").
		Action(func(ctx context.Context, c *Command) error {
			return expectedErr
		})

	err := cmd.ExecuteWithArgs([]string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != expectedErr {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
}

// TestExecuteMissingRequiredArg tests missing required argument error
func TestExecuteMissingRequiredArg(t *testing.T) {
	cmd := Root("test").
		Arg("name", "Name", true).
		Action(func(ctx context.Context, c *Command, name string) error {
			return nil
		})

	err := cmd.ExecuteWithArgs([]string{})
	if err == nil {
		t.Fatal("expected error for missing required argument")
	}
	if _, ok := err.(*ArgumentError); !ok {
		t.Errorf("expected ArgumentError, got %T", err)
	}
}

// TestExecuteInvalidArgumentType tests invalid argument type error
func TestExecuteInvalidArgumentType(t *testing.T) {
	cmd := Root("test").
		Arg("count", "Count", true).
		Action(func(ctx context.Context, c *Command, count int) error {
			return nil
		})

	err := cmd.ExecuteWithArgs([]string{"not-a-number"})
	if err == nil {
		t.Fatal("expected error for invalid argument type")
	}
	if _, ok := err.(*ArgumentError); !ok {
		t.Errorf("expected ArgumentError, got %T", err)
	}
}

// TestExecuteCommandNotFound tests command not found error
func TestExecuteCommandNotFound(t *testing.T) {
	root := Root("app")
	deploy := Cmd("deploy").Action(func(ctx context.Context, c *Command) error {
		return nil
	})
	root.AddCommand(deploy)

	err := root.ExecuteWithArgs([]string{"nonexistent"})
	if err == nil {
		t.Fatal("expected error for nonexistent command")
	}
	if _, ok := err.(*CommandNotFoundError); !ok {
		t.Errorf("expected CommandNotFoundError, got %T", err)
	}
}

// TestExecuteFlagsBeforeSubcommand tests flags before subcommand
func TestExecuteFlagsBeforeSubcommand(t *testing.T) {
	var verbose bool
	var childExecuted bool

	root := Root("app").
		Flag(&verbose, "verbose", "v", false, "Verbose")

	child := Cmd("deploy").
		Action(func(ctx context.Context, c *Command) error {
			childExecuted = true
			return nil
		})

	root.AddCommand(child)

	err := root.ExecuteWithArgs([]string{"--verbose", "deploy"})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}
	if !childExecuted {
		t.Error("child should execute")
	}

	flag := root.flags.GetFlag("verbose")
	if flag != nil {
		if val, ok := flag.GetValue().(bool); !ok || !val {
			t.Error("verbose flag should be true")
		}
	}
} // TestExecuteFlagsAfterSubcommand tests flags after subcommand
func TestExecuteFlagsAfterSubcommand(t *testing.T) {
	var timeout int
	var childExecuted bool

	root := Root("app")

	child := Cmd("deploy").
		Flag(&timeout, "timeout", "t", 30, "Timeout").
		Action(func(ctx context.Context, c *Command) error {
			childExecuted = true
			return nil
		})

	root.AddCommand(child)

	err := root.ExecuteWithArgs([]string{"deploy", "--timeout=60"})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}
	if !childExecuted {
		t.Error("child should execute")
	}

	flag := child.flags.GetFlag("timeout")
	if flag != nil {
		if val, ok := flag.GetValue().(int); !ok || val != 60 {
			t.Errorf("timeout should be 60, got %v", flag.GetValue())
		}
	}
}
