package cli

import (
	"context"
	"testing"
)

// TestCommandCreation tests basic command creation
func TestCommandCreation(t *testing.T) {
	tests := []struct {
		name         string
		createCmd    func() *Command
		expectName   string
		expectDesc   string
		expectHidden bool
	}{
		{
			name: "root command",
			createCmd: func() *Command {
				return Root("myapp").Description("My application")
			},
			expectName:   "myapp",
			expectDesc:   "My application",
			expectHidden: false,
		},
		{
			name: "regular command",
			createCmd: func() *Command {
				return Cmd("deploy").Description("Deploy app").Hidden()
			},
			expectName:   "deploy",
			expectDesc:   "Deploy app",
			expectHidden: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.createCmd()
			if cmd.GetName() != tt.expectName {
				t.Errorf("expected name %q, got %q", tt.expectName, cmd.GetName())
			}
			if cmd.GetDescription() != tt.expectDesc {
				t.Errorf("expected description %q, got %q", tt.expectDesc, cmd.GetDescription())
			}
			if cmd.IsHidden() != tt.expectHidden {
				t.Errorf("expected hidden %v, got %v", tt.expectHidden, cmd.IsHidden())
			}
		})
	}
}

// TestCommandHierarchy tests parent-child relationships
func TestCommandHierarchy(t *testing.T) {
	root := Root("app")
	child := Cmd("child")
	grandchild := Cmd("grandchild")

	root.AddCommand(child)
	child.AddCommand(grandchild)

	if child.GetParent() != root {
		t.Error("child parent should be root")
	}
	if grandchild.GetParent() != child {
		t.Error("grandchild parent should be child")
	}

	rootCommands := root.GetCommands()
	if len(rootCommands) != 1 {
		t.Error("root should have 1 child command")
	}
	if rootCommands["child"] != child {
		t.Error("root's child should be the child command")
	}

	childCommands := child.GetCommands()
	if len(childCommands) != 1 {
		t.Error("child should have 1 grandchild command")
	}
	if childCommands["grandchild"] != grandchild {
		t.Error("child's grandchild should be the grandchild command")
	}
}

// TestCommandArguments tests argument configuration
func TestCommandArguments(t *testing.T) {
	cmd := Cmd("test").
		Arg("name", "User name", true).
		Arg("age", "User age", false)

	args := cmd.GetArgs()
	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(args))
	}

	if args[0].Name != "name" || !args[0].Required {
		t.Error("first arg should be required 'name'")
	}
	if args[1].Name != "age" || args[1].Required {
		t.Error("second arg should be optional 'age'")
	}
}

// TestCommandAction tests action function setting
func TestCommandAction(t *testing.T) {
	executed := false
	cmd := Root("test").
		Action(func(ctx context.Context, cmd *Command) error {
			executed = true
			return nil
		})

	if err := cmd.ExecuteWithArgs([]string{}); err != nil {
		t.Fatalf("execution failed: %v", err)
	}

	if !executed {
		t.Error("action was not executed")
	}
}

// TestCommandHiddenSubcommands tests that hidden commands are excluded from GetCommands
func TestCommandHiddenSubcommands(t *testing.T) {
	root := Root("app")
	visible := Cmd("visible")
	hidden := Cmd("hidden").Hidden()

	root.AddCommand(visible)
	root.AddCommand(hidden)

	commands := root.GetCommands()

	// GetCommands should return all commands, including hidden ones
	// (filtering happens at display time, not retrieval)
	if len(commands) != 2 {
		t.Errorf("expected 2 commands, got %d", len(commands))
	}
}

// TestMultipleActionsOnSameCommand tests that only the last action is kept
func TestMultipleActionsOnSameCommand(t *testing.T) {
	counter := 0
	cmd := Root("test").
		Action(func(ctx context.Context, cmd *Command) error {
			counter = 1
			return nil
		}).
		Action(func(ctx context.Context, cmd *Command) error {
			counter = 2
			return nil
		})

	if err := cmd.ExecuteWithArgs([]string{}); err != nil {
		t.Fatalf("execution failed: %v", err)
	}

	if counter != 2 {
		t.Errorf("expected counter=2 (last action), got %d", counter)
	}
}

// TestCommandChaining tests fluent API chaining
func TestCommandChaining(t *testing.T) {
	cmd := Root("app").
		Description("Test app").
		Arg("name", "Name arg", true).
		Hidden().
		Action(func(ctx context.Context, c *Command) error {
			return nil
		})

	if cmd.GetName() != "app" {
		t.Error("chaining broke name")
	}
	if cmd.GetDescription() != "Test app" {
		t.Error("chaining broke description")
	}
	if !cmd.IsHidden() {
		t.Error("chaining broke hidden")
	}
	if len(cmd.GetArgs()) != 1 {
		t.Error("chaining broke args")
	}
}
