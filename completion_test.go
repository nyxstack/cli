package cli

import (
	"strings"
	"testing"
)

// TestCompletionBash tests bash completion functionality
func TestCompletionBash(t *testing.T) {
	root := Root("myapp").
		Description("Test app")

	deploy := Cmd("deploy").
		Description("Deploy command")

	database := Cmd("database").
		Description("Database command")

	root.AddCommand(deploy)
	root.AddCommand(database)

	bash := &BashCompletion{}

	t.Run("GetCompletions returns subcommands", func(t *testing.T) {
		completions := bash.GetCompletions(root, nil)
		if len(completions) < 2 {
			t.Errorf("expected at least 2 completions, got %d", len(completions))
		}

		hasDeploy := false
		hasDatabase := false
		for _, c := range completions {
			if c == "deploy" {
				hasDeploy = true
			}
			if c == "database" {
				hasDatabase = true
			}
		}

		if !hasDeploy {
			t.Error("completions should include 'deploy'")
		}
		if !hasDatabase {
			t.Error("completions should include 'database'")
		}
	})

	t.Run("GenerateScript returns valid bash script", func(t *testing.T) {
		script := bash.GenerateScript(root)
		if script == "" {
			t.Error("script should not be empty")
		}
		if !strings.Contains(script, "myapp") {
			t.Error("script should contain command name")
		}
		if !strings.Contains(script, "completion") {
			t.Error("script should contain completion logic")
		}
	})
}

// TestCompletionZsh tests zsh completion functionality
func TestCompletionZsh(t *testing.T) {
	root := Root("myapp")
	zsh := &ZshCompletion{}

	script := zsh.GenerateScript(root)
	if script == "" {
		t.Error("zsh script should not be empty")
	}
	if !strings.Contains(script, "myapp") {
		t.Error("zsh script should contain command name")
	}
}

// TestCompletionFish tests fish completion functionality
func TestCompletionFish(t *testing.T) {
	root := Root("myapp")
	fish := &FishCompletion{}

	script := fish.GenerateScript(root)
	if script == "" {
		t.Error("fish script should not be empty")
	}
	if !strings.Contains(script, "myapp") {
		t.Error("fish script should contain command name")
	}
}

// TestCompletionPowerShell tests PowerShell completion functionality
func TestCompletionPowerShell(t *testing.T) {
	root := Root("myapp")
	ps := &PowerShellCompletion{}

	script := ps.GenerateScript(root)
	if script == "" {
		t.Error("PowerShell script should not be empty")
	}
	if !strings.Contains(script, "myapp") {
		t.Error("PowerShell script should contain command name")
	}
}

// TestCompletionWithFlags tests that completions include flags
func TestCompletionWithFlags(t *testing.T) {
	var verbose bool
	var timeout int

	root := Root("myapp").
		Flag(&verbose, "verbose", "v", false, "Verbose").
		Flag(&timeout, "timeout", "t", 30, "Timeout")

	bash := &BashCompletion{}
	completions := bash.GetCompletions(root, nil)

	hasVerbose := false
	hasTimeout := false
	for _, c := range completions {
		if c == "--verbose" || c == "-v" {
			hasVerbose = true
		}
		if c == "--timeout" || c == "-t" {
			hasTimeout = true
		}
	}

	if !hasVerbose {
		t.Error("completions should include verbose flag")
	}
	if !hasTimeout {
		t.Error("completions should include timeout flag")
	}
}

// TestCompletionHiddenCommandsExcluded tests that hidden commands are excluded
func TestCompletionHiddenCommandsExcluded(t *testing.T) {
	root := Root("myapp")

	visible := Cmd("deploy").
		Description("Deploy command")

	hidden := Cmd("debug").
		Description("Debug command").
		Hidden()

	root.AddCommand(visible)
	root.AddCommand(hidden)

	bash := &BashCompletion{}
	completions := bash.GetCompletions(root, nil)

	for _, c := range completions {
		if c == "debug" {
			t.Error("hidden command 'debug' should not be in completions")
		}
	}
}

// TestCompletionRegistersSubcommands tests that completion registers for all subcommands
func TestCompletionRegistersSubcommands(t *testing.T) {
	root := Root("myapp")

	deploy := Cmd("deploy")
	database := Cmd("database")

	root.AddCommand(deploy)
	root.AddCommand(database)

	bash := &BashCompletion{}
	bash.Register(root)

	// Check that __bashcomplete was added to root
	if _, exists := root.subcommands["__bashcomplete"]; !exists {
		t.Error("root should have __bashcomplete command")
	}

	// Check that __bashcomplete was added to subcommands
	if _, exists := deploy.subcommands["__bashcomplete"]; !exists {
		t.Error("deploy should have __bashcomplete command")
	}
	if _, exists := database.subcommands["__bashcomplete"]; !exists {
		t.Error("database should have __bashcomplete command")
	}
}

// TestAddCompletion tests the AddCompletion helper
func TestAddCompletion(t *testing.T) {
	root := Root("myapp")
	deploy := Cmd("deploy")
	root.AddCommand(deploy)

	// Add all completion types
	AddCompletion(root)

	// Check that all completion commands were registered
	completionCommands := []string{
		"__bashcomplete",
		"__zshcomplete",
		"__fishcomplete",
		"__powershellcomplete",
	}

	for _, cmdName := range completionCommands {
		if _, exists := root.subcommands[cmdName]; !exists {
			t.Errorf("root should have %s command", cmdName)
		}
		if _, exists := deploy.subcommands[cmdName]; !exists {
			t.Errorf("deploy should have %s command", cmdName)
		}
	}
}
