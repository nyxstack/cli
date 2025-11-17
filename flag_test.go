package cli

import (
	"testing"
	"time"
)

// TestFlagCreation tests basic flag creation
func TestFlagCreation(t *testing.T) {
	var verbose bool
	var count int
	var timeout time.Duration

	cmd := Root("test").
		Flag(&verbose, "verbose", "v", false, "Enable verbose").
		Flag(&count, "count", "c", 5, "Number of items").
		Flag(&timeout, "timeout", "t", 30*time.Second, "Timeout duration")

	flags := cmd.flags.GetAll()
	if len(flags) != 3 {
		t.Fatalf("expected 3 flags, got %d", len(flags))
	}

	// Check verbose flag
	vFlag := cmd.flags.GetFlag("verbose")
	if vFlag == nil {
		t.Fatal("verbose flag not found")
	}
	if vFlag.GetType() != "bool" {
		t.Errorf("verbose flag type: expected bool, got %s", vFlag.GetType())
	}
	if vFlag.PrimaryName() != "verbose" {
		t.Errorf("verbose primary name: expected 'verbose', got %s", vFlag.PrimaryName())
	}
	if vFlag.ShortName() != "v" {
		t.Errorf("verbose short name: expected 'v', got %s", vFlag.ShortName())
	}

	// Check count flag
	cFlag := cmd.flags.GetFlag("count")
	if cFlag == nil {
		t.Fatal("count flag not found")
	}
	if cFlag.GetType() != "int" {
		t.Errorf("count flag type: expected int, got %s", cFlag.GetType())
	}

	// Check timeout flag
	tFlag := cmd.flags.GetFlag("timeout")
	if tFlag == nil {
		t.Fatal("timeout flag not found")
	}
	if tFlag.GetType() != "duration" {
		t.Errorf("timeout flag type: expected duration, got %s", tFlag.GetType())
	}
}

// TestFlagParsing tests flag value parsing
func TestFlagParsing(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		checkFunc func(*testing.T, *Command)
	}{
		{
			name: "boolean flag",
			args: []string{"--verbose"},
			checkFunc: func(t *testing.T, cmd *Command) {
				var verbose bool
				cmd.Flag(&verbose, "verbose", "", false, "")
				cmd.ExecuteWithArgs([]string{"--verbose"})
				flag := cmd.flags.GetFlag("verbose")
				if val, _ := flag.GetValue().(bool); !val {
					t.Error("verbose should be true")
				}
			},
		},
		{
			name: "int flag",
			args: []string{"--count=42"},
			checkFunc: func(t *testing.T, cmd *Command) {
				var count int
				cmd.Flag(&count, "count", "", 0, "")
				cmd.ExecuteWithArgs([]string{"--count=42"})
				flag := cmd.flags.GetFlag("count")
				if val, _ := flag.GetValue().(int); val != 42 {
					t.Errorf("count should be 42, got %d", val)
				}
			},
		},
		{
			name: "short flag",
			args: []string{"-v"},
			checkFunc: func(t *testing.T, cmd *Command) {
				var verbose bool
				cmd.Flag(&verbose, "verbose", "v", false, "")
				cmd.ExecuteWithArgs([]string{"-v"})
				flag := cmd.flags.GetFlag("verbose")
				if val, _ := flag.GetValue().(bool); !val {
					t.Error("verbose should be true")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Root("test")
			tt.checkFunc(t, cmd)
		})
	}
}

// TestFlagRequired tests required flag validation
func TestFlagRequired(t *testing.T) {
	var apiKey string
	cmd := Root("test").
		FlagRequired(&apiKey, "api-key", "", "", "API key").
		Action(func(ctx interface{}, c *Command) error {
			return nil
		})

	// Should fail without required flag
	err := cmd.ExecuteWithArgs([]string{})
	if err == nil {
		t.Error("expected error for missing required flag")
	}
	if _, ok := err.(*FlagError); !ok {
		t.Errorf("expected FlagError, got %T", err)
	}
}

// TestFlagHidden tests hidden flags
func TestFlagHidden(t *testing.T) {
	var debug bool
	cmd := Root("test").
		FlagHidden(&debug, "debug", "", false, "Debug mode")

	flag := cmd.flags.GetFlag("debug")
	if flag == nil {
		t.Fatal("debug flag not found")
	}
	if !flag.IsHidden() {
		t.Error("debug flag should be hidden")
	}
}

// TestArrayFlags tests array flag handling
func TestArrayFlags(t *testing.T) {
	var tags []string
	cmd := Root("test").
		Flag(&tags, "tag", "t", []string{}, "Tags").
		Action(func(ctx interface{}, c *Command, args ...string) error {
			return nil
		})

	err := cmd.ExecuteWithArgs([]string{"--tag=v1", "--tag=v2", "--tag=v3"})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}

	flag := cmd.flags.GetFlag("tag")
	values := flag.GetValue().([]string)
	if len(values) != 3 {
		t.Fatalf("expected 3 tag values, got %d", len(values))
	}
	if values[0] != "v1" || values[1] != "v2" || values[2] != "v3" {
		t.Errorf("unexpected tag values: %v", values)
	}
}

// TestFlagByShortName tests finding flags by short name
func TestFlagByShortName(t *testing.T) {
	var verbose bool
	cmd := Root("test").
		Flag(&verbose, "verbose", "v", false, "Verbose output")

	flag := cmd.flags.GetFlag("v")
	if flag == nil {
		t.Fatal("should find flag by short name")
	}
	if flag.PrimaryName() != "verbose" {
		t.Error("short name should map to primary flag")
	}
}

// TestFlagUsage tests flag usage strings
func TestFlagUsage(t *testing.T) {
	var verbose bool
	cmd := Root("test").
		Flag(&verbose, "verbose", "v", false, "Enable verbose output")

	flag := cmd.flags.GetFlag("verbose")
	if flag.GetUsage() != "Enable verbose output" {
		t.Errorf("unexpected usage: %s", flag.GetUsage())
	}
}

// TestFlagDefault tests flag default values
func TestFlagDefault(t *testing.T) {
	var count int
	cmd := Root("test").
		Flag(&count, "count", "", 42, "Item count")

	flag := cmd.flags.GetFlag("count")
	defaultVal := flag.GetDefault().(int)
	if defaultVal != 42 {
		t.Errorf("expected default 42, got %d", defaultVal)
	}
}

// TestFlagHasName tests the HasName method
func TestFlagHasName(t *testing.T) {
	var verbose bool
	cmd := Root("test").
		Flag(&verbose, "verbose", "v", false, "Verbose")

	flag := cmd.flags.GetFlag("verbose")
	if !flag.HasName("verbose") {
		t.Error("should have name 'verbose'")
	}
	if !flag.HasName("v") {
		t.Error("should have name 'v'")
	}
	if flag.HasName("other") {
		t.Error("should not have name 'other'")
	}
}
