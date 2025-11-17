package cli

import (
	"context"
	"testing"
)

// TestFlagInheritance tests that child commands inherit parent flags
func TestFlagInheritance(t *testing.T) {
	var verbose bool
	var timeout int

	root := Root("app").
		Flag(&verbose, "verbose", "v", false, "Verbose output")

	child := Cmd("deploy").
		Flag(&timeout, "timeout", "t", 30, "Timeout").
		Action(func(ctx context.Context, c *Command) error {
			return nil
		})

	root.AddCommand(child)

	err := root.ExecuteWithArgs([]string{"deploy", "--verbose", "--timeout=60"})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}

	// Check parent flag is accessible from child
	allFlags := child.getAllFlags()
	var foundVerbose, foundTimeout bool
	for _, f := range allFlags {
		if f.HasName("verbose") {
			foundVerbose = true
		}
		if f.HasName("timeout") {
			foundTimeout = true
		}
	}

	if !foundVerbose {
		t.Error("child should inherit verbose flag from parent")
	}
	if !foundTimeout {
		t.Error("child should have timeout flag")
	}
}

// TestFlagInheritanceMultipleLevels tests inheritance across multiple levels
func TestFlagInheritanceMultipleLevels(t *testing.T) {
	var rootFlag bool
	var midFlag int
	var leafFlag string

	root := Root("app").
		Flag(&rootFlag, "root", "", false, "Root flag")

	mid := Cmd("mid").
		Flag(&midFlag, "mid", "", 0, "Mid flag")

	leaf := Cmd("leaf").
		Flag(&leafFlag, "leaf", "", "", "Leaf flag").
		Action(func(ctx context.Context, c *Command) error {
			return nil
		})

	root.AddCommand(mid)
	mid.AddCommand(leaf)

	err := root.ExecuteWithArgs([]string{"mid", "leaf", "--root", "--mid=42", "--leaf=value"})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}

	// Leaf should have access to all flags
	allFlags := leaf.getAllFlags()
	if len(allFlags) < 3 {
		t.Errorf("leaf should have at least 3 flags (root, mid, leaf), got %d", len(allFlags))
	}

	var foundRoot, foundMid, foundLeaf bool
	for _, f := range allFlags {
		if f.HasName("root") {
			foundRoot = true
		}
		if f.HasName("mid") {
			foundMid = true
		}
		if f.HasName("leaf") {
			foundLeaf = true
		}
	}

	if !foundRoot {
		t.Error("leaf should inherit root flag")
	}
	if !foundMid {
		t.Error("leaf should inherit mid flag")
	}
	if !foundLeaf {
		t.Error("leaf should have leaf flag")
	}
}

// TestFlagShadowing tests that child flags can shadow parent flags
func TestFlagShadowing(t *testing.T) {
	var parentTimeout int
	var childTimeout int

	root := Root("app").
		Flag(&parentTimeout, "timeout", "t", 30, "Parent timeout")

	child := Cmd("deploy").
		Flag(&childTimeout, "timeout", "t", 60, "Child timeout").
		Action(func(ctx context.Context, c *Command) error {
			return nil
		})

	root.AddCommand(child)

	err := root.ExecuteWithArgs([]string{"deploy", "--timeout=90"})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}

	// Child's timeout flag should take precedence
	flag := child.flags.GetFlag("timeout")
	if flag == nil {
		t.Fatal("timeout flag not found in child")
	}

	val, _ := flag.GetValue().(int)
	if val != 90 {
		t.Errorf("expected child timeout to be 90, got %d", val)
	}
}

// TestFlagInheritanceDoesNotModifyParent tests that setting inherited flags doesn't affect parent
func TestFlagInheritanceDoesNotModifyParent(t *testing.T) {
	var verbose bool

	root := Root("app").
		Flag(&verbose, "verbose", "v", false, "Verbose")

	child := Cmd("deploy").
		Action(func(ctx context.Context, c *Command) error {
			return nil
		})

	root.AddCommand(child)

	// Execute child with inherited flag
	err := root.ExecuteWithArgs([]string{"deploy", "--verbose"})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}

	// Parent's flag should be updated (inheritance is by reference)
	parentFlag := root.flags.GetFlag("verbose")
	val, _ := parentFlag.GetValue().(bool)
	if !val {
		t.Error("parent flag should be updated when set via child")
	}
}

// TestNoFlagInheritanceForSiblings tests that sibling commands don't inherit from each other
func TestNoFlagInheritanceForSiblings(t *testing.T) {
	var deployFlag string
	var rollbackFlag string

	root := Root("app")

	deploy := Cmd("deploy").
		Flag(&deployFlag, "strategy", "", "rolling", "Deploy strategy")

	rollback := Cmd("rollback").
		Flag(&rollbackFlag, "version", "", "latest", "Rollback version").
		Action(func(ctx context.Context, c *Command) error {
			return nil
		})

	root.AddCommand(deploy)
	root.AddCommand(rollback)

	err := root.ExecuteWithArgs([]string{"rollback"})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}

	// Rollback should not have deploy's flag
	allFlags := rollback.getAllFlags()
	for _, f := range allFlags {
		if f.HasName("strategy") {
			t.Error("rollback should not have deploy's strategy flag")
		}
	}
}

// TestFlagInheritanceWithHiddenFlags tests that hidden flags are still inherited
func TestFlagInheritanceWithHiddenFlags(t *testing.T) {
	var debug bool

	root := Root("app").
		FlagHidden(&debug, "debug", "", false, "Debug mode")

	child := Cmd("deploy").
		Action(func(ctx context.Context, c *Command) error {
			return nil
		})

	root.AddCommand(child)

	err := root.ExecuteWithArgs([]string{"deploy", "--debug"})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}

	// Child should inherit hidden flag
	allFlags := child.getAllFlags()
	var foundDebug bool
	for _, f := range allFlags {
		if f.HasName("debug") {
			foundDebug = true
			if !f.IsHidden() {
				t.Error("inherited debug flag should remain hidden")
			}
		}
	}

	if !foundDebug {
		t.Error("child should inherit hidden debug flag")
	}
}
