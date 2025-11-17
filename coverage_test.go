package cli

import (
	"context"
	"testing"
	"time"
)

// TestCommandVisibility tests Show/Hidden methods
func TestCommandVisibility(t *testing.T) {
	cmd := Cmd("test").Hidden()
	if !cmd.IsHidden() {
		t.Error("command should be hidden")
	}

	cmd.Show()
	if cmd.IsHidden() {
		t.Error("command should be visible after Show()")
	}
}

// TestHelpConfiguration tests help flag customization
func TestHelpConfiguration(t *testing.T) {
	cmd := Root("test")

	// Test IsHelpEnabled (default true)
	if !cmd.IsHelpEnabled() {
		t.Error("help should be enabled by default")
	}

	// Test DisableHelp
	cmd.DisableHelp()
	if cmd.IsHelpEnabled() {
		t.Error("help should be disabled after DisableHelp()")
	}

	// Test EnableHelp
	cmd.EnableHelp()
	if !cmd.IsHelpEnabled() {
		t.Error("help should be enabled after EnableHelp()")
	}

	// Test SetHelpFlag
	cmd.SetHelpFlag("assist", "a")
	if cmd.helpFlag != "assist" || cmd.helpShort != "a" {
		t.Errorf("expected helpFlag='assist', helpShort='a', got helpFlag='%s', helpShort='%s'", cmd.helpFlag, cmd.helpShort)
	}
}

// TestFlagGetters tests flag getter methods
func TestFlagGetters(t *testing.T) {
	var port int
	var verbose bool

	cmd := Root("test").
		Flag(&port, "port", "p", 8080, "Port number").
		Flag(&verbose, "verbose", "v", false, "Verbose output")

	flags := cmd.flags.GetFlags()
	if len(flags) != 2 {
		t.Fatalf("expected 2 flags, got %d", len(flags))
	}

	// Test GetNames
	portFlag := cmd.flags.GetFlag("port")
	names := portFlag.GetNames()
	if len(names) != 2 || names[0] != "port" || names[1] != "p" {
		t.Errorf("expected names [port, p], got %v", names)
	}

	// Test GetType
	flagType := portFlag.GetType()
	if flagType != "int" {
		t.Errorf("expected type 'int', got %q", flagType)
	}

	// Test PrimaryName
	if portFlag.PrimaryName() != "port" {
		t.Errorf("expected primary name 'port', got %q", portFlag.PrimaryName())
	}

	// Test ShortName
	if portFlag.ShortName() != "p" {
		t.Errorf("expected short name 'p', got %q", portFlag.ShortName())
	}

	// Test GetDefault
	if portFlag.GetDefault() != 8080 {
		t.Errorf("expected default 8080, got %v", portFlag.GetDefault())
	}

	// Test GetUsage
	if portFlag.GetUsage() != "Port number" {
		t.Errorf("expected usage 'Port number', got %q", portFlag.GetUsage())
	}

	// Test GetValue before setting
	val := portFlag.GetValue()
	if val != 8080 {
		t.Errorf("expected value 8080, got %v", val)
	}
}

// TestFlagShortNameHandling tests flags with no short name
func TestFlagShortNameHandling(t *testing.T) {
	var config string
	cmd := Root("test").Flag(&config, "config", "", "default.yaml", "Config file")

	flag := cmd.flags.GetFlag("config")
	if flag.ShortName() != "" {
		t.Errorf("expected empty short name, got %q", flag.ShortName())
	}

	if flag.PrimaryName() != "config" {
		t.Errorf("expected primary name 'config', got %q", flag.PrimaryName())
	}
}

// TestCompletionShells tests all shell completion implementations
func TestCompletionShells(t *testing.T) {
	root := Root("test").Description("Test app")

	deploy := Cmd("deploy").Description("Deploy")
	root.AddCommand(deploy)

	tests := []struct {
		name       string
		completion ShellCompletion
	}{
		{"bash", &BashCompletion{}},
		{"zsh", &ZshCompletion{}},
		{"fish", &FishCompletion{}},
		{"powershell", &PowerShellCompletion{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test GetCompletions with empty args
			words := tt.completion.GetCompletions(root, []string{})
			found := false
			for _, word := range words {
				if word == "deploy" {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("%s completion should include 'deploy' command", tt.name)
			}

			// Test GenerateScript (just verify it doesn't panic and returns non-empty)
			script := tt.completion.GenerateScript(root)
			if len(script) == 0 {
				t.Errorf("%s GenerateScript returned empty script", tt.name)
			}

			// Test Register (verify it doesn't panic)
			tt.completion.Register(root)
		})
	}
}

// TestStructBasedFlags tests Flags() method with struct tags
func TestStructBasedFlags(t *testing.T) {
	type Config struct {
		Host    string        `cli:"host,h" default:"localhost" usage:"Server host"`
		Port    int           `cli:"port,p" default:"8080" usage:"Server port"`
		Verbose bool          `cli:"verbose,v" default:"false" usage:"Verbose output"`
		Timeout time.Duration `cli:"timeout,t" default:"30s" usage:"Request timeout"`
	}

	var config Config
	cmd := Root("test").Flags(&config)

	// Verify flags were added
	hostFlag := cmd.flags.GetFlag("host")
	if hostFlag == nil {
		t.Fatal("host flag not found")
	}
	if hostFlag.PrimaryName() != "host" {
		t.Errorf("expected primary name 'host', got %q", hostFlag.PrimaryName())
	}
	if hostFlag.ShortName() != "h" {
		t.Errorf("expected short name 'h', got %q", hostFlag.ShortName())
	}

	portFlag := cmd.flags.GetFlag("port")
	if portFlag == nil {
		t.Fatal("port flag not found")
	}

	verboseFlag := cmd.flags.GetFlag("verbose")
	if verboseFlag == nil {
		t.Fatal("verbose flag not found")
	}

	timeoutFlag := cmd.flags.GetFlag("timeout")
	if timeoutFlag == nil {
		t.Fatal("timeout flag not found")
	}

	// Test execution with struct flags
	err := cmd.ExecuteWithArgs([]string{"--host=example.com", "--port=9000", "--verbose", "--timeout=1m"})
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}

	if config.Host != "example.com" {
		t.Errorf("expected host 'example.com', got %q", config.Host)
	}
	if config.Port != 9000 {
		t.Errorf("expected port 9000, got %d", config.Port)
	}
	if !config.Verbose {
		t.Error("expected verbose to be true")
	}
	if config.Timeout != time.Minute {
		t.Errorf("expected timeout 1m, got %v", config.Timeout)
	}
}

// TestFlagTypeConversions tests all type conversions in setValue
func TestFlagTypeConversions(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() (*Command, interface{})
		args     []string
		validate func(t *testing.T, val interface{})
	}{
		{
			name: "uint conversion",
			setup: func() (*Command, interface{}) {
				var count uint
				cmd := Root("test").
					Flag(&count, "count", "c", uint(10), "Count").
					Action(func(ctx context.Context, cmd *Command) error { return nil })
				return cmd, &count
			},
			args: []string{"--count=42"},
			validate: func(t *testing.T, val interface{}) {
				if *val.(*uint) != 42 {
					t.Errorf("expected 42, got %d", *val.(*uint))
				}
			},
		},
		{
			name: "float32 conversion",
			setup: func() (*Command, interface{}) {
				var ratio float32
				cmd := Root("test").
					Flag(&ratio, "ratio", "r", float32(1.5), "Ratio").
					Action(func(ctx context.Context, cmd *Command) error { return nil })
				return cmd, &ratio
			},
			args: []string{"--ratio=3.14"},
			validate: func(t *testing.T, val interface{}) {
				expected := float32(3.14)
				actual := *val.(*float32)
				if actual < expected-0.01 || actual > expected+0.01 {
					t.Errorf("expected ~3.14, got %f", actual)
				}
			},
		},
		{
			name: "float64 conversion",
			setup: func() (*Command, interface{}) {
				var ratio float64
				cmd := Root("test").
					Flag(&ratio, "ratio", "r", 1.5, "Ratio").
					Action(func(ctx context.Context, cmd *Command) error { return nil })
				return cmd, &ratio
			},
			args: []string{"--ratio=2.71828"},
			validate: func(t *testing.T, val interface{}) {
				if *val.(*float64) != 2.71828 {
					t.Errorf("expected 2.71828, got %f", *val.(*float64))
				}
			},
		},
		{
			name: "int8 conversion",
			setup: func() (*Command, interface{}) {
				var level int8
				cmd := Root("test").
					Flag(&level, "level", "l", int8(1), "Level").
					Action(func(ctx context.Context, cmd *Command) error { return nil })
				return cmd, &level
			},
			args: []string{"--level=127"},
			validate: func(t *testing.T, val interface{}) {
				if *val.(*int8) != 127 {
					t.Errorf("expected 127, got %d", *val.(*int8))
				}
			},
		},
		{
			name: "int16 conversion",
			setup: func() (*Command, interface{}) {
				var port int16
				cmd := Root("test").
					Flag(&port, "port", "p", int16(8080), "Port").
					Action(func(ctx context.Context, cmd *Command) error { return nil })
				return cmd, &port
			},
			args: []string{"--port=32767"},
			validate: func(t *testing.T, val interface{}) {
				if *val.(*int16) != 32767 {
					t.Errorf("expected 32767, got %d", *val.(*int16))
				}
			},
		},
		{
			name: "int32 conversion",
			setup: func() (*Command, interface{}) {
				var count int32
				cmd := Root("test").
					Flag(&count, "count", "c", int32(100), "Count").
					Action(func(ctx context.Context, cmd *Command) error { return nil })
				return cmd, &count
			},
			args: []string{"--count=2147483647"},
			validate: func(t *testing.T, val interface{}) {
				if *val.(*int32) != 2147483647 {
					t.Errorf("expected 2147483647, got %d", *val.(*int32))
				}
			},
		},
		{
			name: "int64 conversion",
			setup: func() (*Command, interface{}) {
				var size int64
				cmd := Root("test").
					Flag(&size, "size", "s", int64(1000), "Size").
					Action(func(ctx context.Context, cmd *Command) error { return nil })
				return cmd, &size
			},
			args: []string{"--size=9223372036854775807"},
			validate: func(t *testing.T, val interface{}) {
				if *val.(*int64) != 9223372036854775807 {
					t.Errorf("expected 9223372036854775807, got %d", *val.(*int64))
				}
			},
		},
		{
			name: "uint8 conversion",
			setup: func() (*Command, interface{}) {
				var level uint8
				cmd := Root("test").
					Flag(&level, "level", "l", uint8(0), "Level").
					Action(func(ctx context.Context, cmd *Command) error { return nil })
				return cmd, &level
			},
			args: []string{"--level=255"},
			validate: func(t *testing.T, val interface{}) {
				if *val.(*uint8) != 255 {
					t.Errorf("expected 255, got %d", *val.(*uint8))
				}
			},
		},
		{
			name: "uint16 conversion",
			setup: func() (*Command, interface{}) {
				var port uint16
				cmd := Root("test").
					Flag(&port, "port", "p", uint16(8080), "Port").
					Action(func(ctx context.Context, cmd *Command) error { return nil })
				return cmd, &port
			},
			args: []string{"--port=65535"},
			validate: func(t *testing.T, val interface{}) {
				if *val.(*uint16) != 65535 {
					t.Errorf("expected 65535, got %d", *val.(*uint16))
				}
			},
		},
		{
			name: "uint32 conversion",
			setup: func() (*Command, interface{}) {
				var count uint32
				cmd := Root("test").
					Flag(&count, "count", "c", uint32(100), "Count").
					Action(func(ctx context.Context, cmd *Command) error { return nil })
				return cmd, &count
			},
			args: []string{"--count=4294967295"},
			validate: func(t *testing.T, val interface{}) {
				if *val.(*uint32) != 4294967295 {
					t.Errorf("expected 4294967295, got %d", *val.(*uint32))
				}
			},
		},
		{
			name: "uint64 conversion",
			setup: func() (*Command, interface{}) {
				var size uint64
				cmd := Root("test").
					Flag(&size, "size", "s", uint64(1000), "Size").
					Action(func(ctx context.Context, cmd *Command) error { return nil })
				return cmd, &size
			},
			args: []string{"--size=18446744073709551615"},
			validate: func(t *testing.T, val interface{}) {
				if *val.(*uint64) != 18446744073709551615 {
					t.Errorf("expected 18446744073709551615, got %d", *val.(*uint64))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, valPtr := tt.setup()
			err := cmd.ExecuteWithArgs(tt.args)
			if err != nil {
				t.Fatalf("execution failed: %v", err)
			}
			tt.validate(t, valPtr)
		})
	}
}

// TestVariadicActionEdgeCases tests edge cases in variadic action handling
func TestVariadicActionEdgeCases(t *testing.T) {
	t.Run("variadic with fixed args", func(t *testing.T) {
		var result []string
		cmd := Root("test").
			Arg("first", "First arg", true).
			Arg("rest", "Rest args", false).
			Action(func(ctx context.Context, c *Command, first string, rest ...string) error {
				result = append([]string{first}, rest...)
				return nil
			})

		err := cmd.ExecuteWithArgs([]string{"one", "two", "three"})
		if err != nil {
			t.Fatalf("execution failed: %v", err)
		}

		expected := []string{"one", "two", "three"}
		if len(result) != len(expected) {
			t.Errorf("expected %d results, got %d", len(expected), len(result))
		}
		for i, v := range expected {
			if result[i] != v {
				t.Errorf("at index %d: expected %q, got %q", i, v, result[i])
			}
		}
	})

	t.Run("variadic with no optional args", func(t *testing.T) {
		var result []string
		cmd := Root("test").
			Arg("first", "First arg", true).
			Action(func(ctx context.Context, c *Command, first string, rest ...string) error {
				result = append([]string{first}, rest...)
				return nil
			})

		err := cmd.ExecuteWithArgs([]string{"only"})
		if err != nil {
			t.Fatalf("execution failed: %v", err)
		}

		if len(result) != 1 || result[0] != "only" {
			t.Errorf("expected [only], got %v", result)
		}
	})
}

// TestInferType tests type inference for struct fields
func TestInferType(t *testing.T) {
	type TestStruct struct {
		StrField    string        `cli:"str-field,s" usage:"String field"`
		IntField    int           `cli:"int-field,i" usage:"Int field"`
		BoolField   bool          `cli:"bool-field,b" usage:"Bool field"`
		FloatField  float64       `cli:"float-field,f" usage:"Float field"`
		DurField    time.Duration `cli:"dur-field,d" usage:"Duration field"`
		Uint64Field uint64        `cli:"uint64-field,u" usage:"Uint64 field"`
	}

	var ts TestStruct
	cmd := Root("test")

	// This will internally call inferType for each field
	cmd.Flags(&ts)

	// Verify flags were created with correct types
	strFlag := cmd.flags.GetFlag("str-field")
	if strFlag == nil {
		t.Fatal("str-field flag not created")
	}
	if strFlag.GetType() != "string" {
		t.Errorf("expected type 'string', got %q", strFlag.GetType())
	}

	intFlag := cmd.flags.GetFlag("int-field")
	if intFlag == nil {
		t.Fatal("int-field flag not created")
	}
	if intFlag.GetType() != "int" {
		t.Errorf("expected type 'int', got %q", intFlag.GetType())
	}

	boolFlag := cmd.flags.GetFlag("bool-field")
	if boolFlag == nil {
		t.Fatal("bool-field flag not created")
	}
	if boolFlag.GetType() != "bool" {
		t.Errorf("expected type 'bool', got %q", boolFlag.GetType())
	}

	floatFlag := cmd.flags.GetFlag("float-field")
	if floatFlag == nil {
		t.Fatal("float-field flag not created")
	}
	if floatFlag.GetType() != "float" {
		t.Errorf("expected type 'float', got %q", floatFlag.GetType())
	}

	durFlag := cmd.flags.GetFlag("dur-field")
	if durFlag == nil {
		t.Fatal("dur-field flag not created")
	}
	if durFlag.GetType() != "duration" {
		t.Errorf("expected type 'duration', got %q", durFlag.GetType())
	}

	uint64Flag := cmd.flags.GetFlag("uint64-field")
	if uint64Flag == nil {
		t.Fatal("uint64-field flag not created")
	}
	if uint64Flag.GetType() != "uint" {
		t.Errorf("expected type 'uint', got %q", uint64Flag.GetType())
	}
}
