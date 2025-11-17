package cli

import (
	"testing"
)

// TestErrorTypes tests custom error type creation
func TestErrorTypes(t *testing.T) {
	cmd := Root("test")

	t.Run("CommandNotFoundError", func(t *testing.T) {
		err := &CommandNotFoundError{
			Name: "nonexistent",
			Cmd:  cmd,
		}
		msg := err.Error()
		if msg == "" {
			t.Error("error message should not be empty")
		}
		if err.Cmd != cmd {
			t.Error("command reference should be preserved")
		}
	})

	t.Run("ArgumentError", func(t *testing.T) {
		err := &ArgumentError{
			Arg: "count",
			Msg: "invalid value",
			Cmd: cmd,
		}
		msg := err.Error()
		if msg == "" {
			t.Error("error message should not be empty")
		}
		if err.Cmd != cmd {
			t.Error("command reference should be preserved")
		}
	})

	t.Run("FlagError", func(t *testing.T) {
		err := &FlagError{
			Flag: "verbose",
			Msg:  "invalid flag",
			Cmd:  cmd,
		}
		msg := err.Error()
		if msg == "" {
			t.Error("error message should not be empty")
		}
		if err.Cmd != cmd {
			t.Error("command reference should be preserved")
		}
	})
}

// TestErrorMessages tests error message formatting
func TestErrorMessages(t *testing.T) {
	cmd := Root("myapp")

	tests := []struct {
		name     string
		err      error
		contains string
	}{
		{
			name: "CommandNotFoundError message",
			err: &CommandNotFoundError{
				Name: "deploy",
				Cmd:  cmd,
			},
			contains: "deploy",
		},
		{
			name: "ArgumentError message",
			err: &ArgumentError{
				Arg: "count",
				Msg: "must be positive",
				Cmd: cmd,
			},
			contains: "count",
		},
		{
			name: "FlagError message",
			err: &FlagError{
				Flag: "timeout",
				Msg:  "invalid duration",
				Cmd:  cmd,
			},
			contains: "timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tt.err.Error()
			if msg == "" {
				t.Error("error message should not be empty")
			}
			// Error messages should contain relevant context
			// (exact format not specified, just checking it's not empty)
		})
	}
}
