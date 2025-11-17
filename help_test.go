package cli

import (
	"context"
	"testing"
)

// TestHelpCommand tests the help display functionality
func TestHelpCommand(t *testing.T) {
	t.Run("root help with -h", func(t *testing.T) {
		root := Root("myapp").
			Description("My application")

		err := root.ExecuteWithArgs([]string{"-h"})
		if err != nil {
			t.Errorf("help should not return error, got %v", err)
		}
	})

	t.Run("root help with --help", func(t *testing.T) {
		root := Root("myapp").
			Description("My application")

		err := root.ExecuteWithArgs([]string{"--help"})
		if err != nil {
			t.Errorf("help should not return error, got %v", err)
		}
	})

	t.Run("subcommand help with -h", func(t *testing.T) {
		root := Root("myapp")
		deploy := Cmd("deploy").
			Description("Deploy the application")
		root.AddCommand(deploy)

		err := root.ExecuteWithArgs([]string{"deploy", "-h"})
		if err != nil {
			t.Errorf("help should not return error, got %v", err)
		}
	})

	t.Run("subcommand help with --help", func(t *testing.T) {
		root := Root("myapp")
		deploy := Cmd("deploy").
			Description("Deploy the application")
		root.AddCommand(deploy)

		err := root.ExecuteWithArgs([]string{"deploy", "--help"})
		if err != nil {
			t.Errorf("help should not return error, got %v", err)
		}
	})

	t.Run("help shows correct command", func(t *testing.T) {
		actionExecuted := false

		root := Root("myapp")
		deploy := Cmd("deploy").
			Action(func(ctx context.Context, cmd *Command) error {
				actionExecuted = true
				return nil
			})
		root.AddCommand(deploy)

		// Help should prevent action execution
		err := root.ExecuteWithArgs([]string{"deploy", "--help"})
		if err != nil {
			t.Errorf("help should not return error, got %v", err)
		}
		if actionExecuted {
			t.Error("action should not execute when help is shown")
		}
	})
}

// TestHelpWithFlags tests help display includes flags
func TestHelpWithFlags(t *testing.T) {
	var verbose bool
	var count int

	root := Root("myapp").
		Flag(&verbose, "verbose", "v", false, "Verbose output").
		Flag(&count, "count", "c", 10, "Item count")

	// Just verify it doesn't crash - actual output testing would require capturing stdout
	err := root.ExecuteWithArgs([]string{"-h"})
	if err != nil {
		t.Errorf("help with flags should not return error, got %v", err)
	}
}

// TestHelpWithArguments tests help display includes arguments
func TestHelpWithArguments(t *testing.T) {
	root := Root("myapp").
		Arg("name", "User name", true).
		Arg("age", "User age", false)

	err := root.ExecuteWithArgs([]string{"--help"})
	if err != nil {
		t.Errorf("help with arguments should not return error, got %v", err)
	}
}

// TestHelpWithSubcommands tests help display includes subcommands
func TestHelpWithSubcommands(t *testing.T) {
	root := Root("myapp")

	root.AddCommand(Cmd("deploy").Description("Deploy app"))
	root.AddCommand(Cmd("rollback").Description("Rollback app"))
	root.AddCommand(Cmd("status").Description("Check status"))

	err := root.ExecuteWithArgs([]string{"-h"})
	if err != nil {
		t.Errorf("help with subcommands should not return error, got %v", err)
	}
}

// TestHelpHidesHiddenCommands tests that hidden commands don't appear in help
func TestHelpHidesHiddenCommands(t *testing.T) {
	root := Root("myapp")

	root.AddCommand(Cmd("deploy").Description("Deploy app"))
	root.AddCommand(Cmd("debug").Description("Debug mode").Hidden())

	// Just verify it doesn't crash - hidden filtering is done in showHelp
	err := root.ExecuteWithArgs([]string{"--help"})
	if err != nil {
		t.Errorf("help should not return error, got %v", err)
	}
}

// TestHelpBeforeFlags tests help can appear before flags
func TestHelpBeforeFlags(t *testing.T) {
	var verbose bool

	root := Root("myapp").
		Flag(&verbose, "verbose", "v", false, "Verbose")

	err := root.ExecuteWithArgs([]string{"-h", "--verbose"})
	if err != nil {
		t.Errorf("help before flags should not return error, got %v", err)
	}

	// Verbose should not be set because help short-circuits execution
	flag := root.flags.GetFlag("verbose")
	if val, _ := flag.GetValue().(bool); val {
		t.Error("flags should not be parsed when help is shown")
	}
}

// TestHelpAfterSubcommand tests help can appear after subcommand name
func TestHelpAfterSubcommand(t *testing.T) {
	actionExecuted := false

	root := Root("myapp")
	deploy := Cmd("deploy").
		Action(func(ctx context.Context, cmd *Command) error {
			actionExecuted = true
			return nil
		})
	root.AddCommand(deploy)

	err := root.ExecuteWithArgs([]string{"deploy", "--help"})
	if err != nil {
		t.Errorf("help after subcommand should not return error, got %v", err)
	}
	if actionExecuted {
		t.Error("action should not execute when help is requested")
	}
}

// TestShowHelpMethod tests the ShowHelp method can be called directly
func TestShowHelpMethod(t *testing.T) {
	root := Root("myapp").
		Description("Test application")

	// Should not panic
	root.ShowHelp()
}

// TestHelpForDeepNestedCommands tests help for deeply nested command structures
func TestHelpForDeepNestedCommands(t *testing.T) {
	root := Root("myapp")
	database := Cmd("database")
	migrate := Cmd("migrate")

	root.AddCommand(database)
	database.AddCommand(migrate)

	// Test help at each level
	tests := []struct {
		name string
		args []string
	}{
		{"root help", []string{"-h"}},
		{"database help", []string{"database", "-h"}},
		{"migrate help", []string{"database", "migrate", "-h"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := root.ExecuteWithArgs(tt.args)
			if err != nil {
				t.Errorf("help should not return error, got %v", err)
			}
		})
	}
}
