package cli

import (
	"context"
	"testing"
)

// TestLifecycleHooksOrder tests that hooks execute in correct order
func TestLifecycleHooksOrder(t *testing.T) {
	order := []string{}

	cmd := Root("test").
		PersistentPreRun(func(ctx context.Context, c *Command) error {
			order = append(order, "persistent-pre")
			return nil
		}).
		PreRun(func(ctx context.Context, c *Command) error {
			order = append(order, "pre")
			return nil
		}).
		Action(func(ctx context.Context, c *Command) error {
			order = append(order, "action")
			return nil
		}).
		PostRun(func(ctx context.Context, c *Command) error {
			order = append(order, "post")
			return nil
		}).
		PersistentPostRun(func(ctx context.Context, c *Command) error {
			order = append(order, "persistent-post")
			return nil
		})

	err := cmd.ExecuteWithArgs([]string{})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}

	expected := []string{"persistent-pre", "pre", "action", "post", "persistent-post"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d hooks, got %d", len(expected), len(order))
	}

	for i, exp := range expected {
		if order[i] != exp {
			t.Errorf("position %d: expected %s, got %s", i, exp, order[i])
		}
	}
}

// TestLifecycleErrorInPreRun tests that error in PreRun stops execution
func TestLifecycleErrorInPreRun(t *testing.T) {
	actionExecuted := false
	postExecuted := false

	cmd := Root("test").
		PreRun(func(ctx context.Context, c *Command) error {
			return &ArgumentError{Arg: "test", Msg: "validation failed", Cmd: c}
		}).
		Action(func(ctx context.Context, c *Command) error {
			actionExecuted = true
			return nil
		}).
		PostRun(func(ctx context.Context, c *Command) error {
			postExecuted = true
			return nil
		})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error from PreRun")
	}
	if actionExecuted {
		t.Error("action should not execute after PreRun error")
	}
	if postExecuted {
		t.Error("PostRun should not execute after PreRun error")
	}
}

// TestLifecycleErrorInAction tests that error in Action still runs PostRun
func TestLifecycleErrorInAction(t *testing.T) {
	postExecuted := false
	persistentPostExecuted := false

	cmd := Root("test").
		Action(func(ctx context.Context, c *Command) error {
			return &ArgumentError{Arg: "test", Msg: "action error", Cmd: c}
		}).
		PostRun(func(ctx context.Context, c *Command) error {
			postExecuted = true
			return nil
		}).
		PersistentPostRun(func(ctx context.Context, c *Command) error {
			persistentPostExecuted = true
			return nil
		})

	err := cmd.ExecuteWithArgs([]string{})
	if err == nil {
		t.Fatal("expected error from Action")
	}
	if !postExecuted {
		t.Error("PostRun should execute even after Action error")
	}
	if !persistentPostExecuted {
		t.Error("PersistentPostRun should execute even after Action error")
	}
}

// TestPersistentHooksInheritance tests persistent hooks run for child commands
func TestPersistentHooksInheritance(t *testing.T) {
	order := []string{}

	root := Root("app").
		PersistentPreRun(func(ctx context.Context, c *Command) error {
			order = append(order, "root-persistent-pre")
			return nil
		}).
		PersistentPostRun(func(ctx context.Context, c *Command) error {
			order = append(order, "root-persistent-post")
			return nil
		})

	child := Cmd("child").
		PreRun(func(ctx context.Context, c *Command) error {
			order = append(order, "child-pre")
			return nil
		}).
		Action(func(ctx context.Context, c *Command) error {
			order = append(order, "child-action")
			return nil
		}).
		PostRun(func(ctx context.Context, c *Command) error {
			order = append(order, "child-post")
			return nil
		})

	root.AddCommand(child)

	err := root.ExecuteWithArgs([]string{"child"})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}

	expected := []string{
		"root-persistent-pre",
		"child-pre",
		"child-action",
		"child-post",
		"root-persistent-post",
	}

	if len(order) != len(expected) {
		t.Fatalf("expected %d hooks, got %d: %v", len(expected), len(order), order)
	}

	for i, exp := range expected {
		if order[i] != exp {
			t.Errorf("position %d: expected %s, got %s", i, exp, order[i])
		}
	}
}

// TestLifecycleNoAction tests execution without Action function
func TestLifecycleNoAction(t *testing.T) {
	preExecuted := false
	postExecuted := false

	cmd := Root("test").
		PreRun(func(ctx context.Context, c *Command) error {
			preExecuted = true
			return nil
		}).
		PostRun(func(ctx context.Context, c *Command) error {
			postExecuted = true
			return nil
		})

	err := cmd.ExecuteWithArgs([]string{})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}
	if !preExecuted {
		t.Error("PreRun should execute")
	}
	if !postExecuted {
		t.Error("PostRun should execute")
	}
}

// TestLifecycleMultiplePersistentLevels tests nested persistent hooks
func TestLifecycleMultiplePersistentLevels(t *testing.T) {
	order := []string{}

	root := Root("app").
		PersistentPreRun(func(ctx context.Context, c *Command) error {
			order = append(order, "root-persistent-pre")
			return nil
		}).
		PersistentPostRun(func(ctx context.Context, c *Command) error {
			order = append(order, "root-persistent-post")
			return nil
		})

	mid := Cmd("mid").
		PersistentPreRun(func(ctx context.Context, c *Command) error {
			order = append(order, "mid-persistent-pre")
			return nil
		}).
		PersistentPostRun(func(ctx context.Context, c *Command) error {
			order = append(order, "mid-persistent-post")
			return nil
		})

	leaf := Cmd("leaf").
		Action(func(ctx context.Context, c *Command) error {
			order = append(order, "leaf-action")
			return nil
		})

	root.AddCommand(mid)
	mid.AddCommand(leaf)

	err := root.ExecuteWithArgs([]string{"mid", "leaf"})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}

	expected := []string{
		"root-persistent-pre",
		"mid-persistent-pre",
		"leaf-action",
		"mid-persistent-post",
		"root-persistent-post",
	}

	if len(order) != len(expected) {
		t.Fatalf("expected %d hooks, got %d: %v", len(expected), len(order), order)
	}

	for i, exp := range expected {
		if order[i] != exp {
			t.Errorf("position %d: expected %s, got %s", i, exp, order[i])
		}
	}
}
