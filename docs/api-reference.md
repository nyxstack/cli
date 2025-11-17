# API Reference

Complete API reference for nyxstack/cli.

## Command

### Creating Commands

```go
func Root(name string) *Command
```
Creates a root command (entry point of your CLI).

```go
func Cmd(name string) *Command
```
Creates a regular subcommand.

### Command Methods

#### Configuration

```go
func (c *Command) Description(desc string) *Command
```
Sets the command description (shown in help).

```go
func (c *Command) Hidden() *Command
```
Hides the command from help output.

```go
func (c *Command) Show() *Command
```
Makes a hidden command visible again.

#### Flags

```go
func (c *Command) Flag(ptr interface{}, name, short string, defaultValue interface{}, usage string) *Command
```
Adds a flag to the command.
- `ptr`: Pointer to variable that will receive the flag value
- `name`: Long flag name (e.g., "verbose")
- `short`: Short flag name (e.g., "v"), use "" for none
- `defaultValue`: Default value if flag not provided
- `usage`: Help text for the flag

```go
func (c *Command) FlagRequired(ptr interface{}, name, short string, defaultValue interface{}, usage string) *Command
```
Adds a required flag. Returns error if not provided.

```go
func (c *Command) FlagHidden(ptr interface{}, name, short string, defaultValue interface{}, usage string) *Command
```
Adds a flag that's hidden from help output.

```go
func (c *Command) Flags(structPtr interface{}) *Command
```
Binds struct fields as flags using struct tags.

Example:
```go
type Config struct {
    Host string `cli:"host,h" default:"localhost" usage:"Server host"`
    Port int    `cli:"port,p" default:"8080" usage:"Server port"`
}
var config Config
cmd.Flags(&config)
```

#### Arguments

```go
func (c *Command) Arg(name, description string, required bool) *Command
```
Adds a positional argument.
- `name`: Argument name (for help text)
- `description`: Argument description
- `required`: Whether argument is required

#### Lifecycle Hooks

```go
func (c *Command) PersistentPreRun(fn ActionFunc) *Command
```
Runs before action, for this command and all children.

```go
func (c *Command) PreRun(fn ActionFunc) *Command
```
Runs before action, only for this command.

```go
func (c *Command) Action(fn ActionFunc) *Command
```
Main command logic.

```go
func (c *Command) PostRun(fn ActionFunc) *Command
```
Runs after action, only for this command.

```go
func (c *Command) PersistentPostRun(fn ActionFunc) *Command
```
Runs after action, for this command and all children.

**ActionFunc signature:**
```go
func(ctx context.Context, cmd *Command, args ...interface{}) error
```

#### Command Hierarchy

```go
func (c *Command) AddCommand(cmd *Command) *Command
```
Adds a subcommand.

#### Execution

```go
func (c *Command) Execute() error
```
Executes the command with os.Args.

```go
func (c *Command) ExecuteContext(ctx context.Context) error
```
Executes with a context (for cancellation/timeout).

```go
func (c *Command) ExecuteWithArgs(args []string) error
```
Executes with custom arguments (useful for testing).

#### Getters

```go
func (c *Command) GetName() string
```
Returns command name.

```go
func (c *Command) GetDescription() string
```
Returns command description.

```go
func (c *Command) GetParent() *Command
```
Returns parent command (nil for root).

```go
func (c *Command) GetCommands() map[string]*Command
```
Returns map of subcommands.

```go
func (c *Command) GetArgs() []Argument
```
Returns list of arguments.

```go
func (c *Command) IsHidden() bool
```
Returns whether command is hidden.

#### Help

```go
func (c *Command) ShowHelp()
```
Displays help for this command.

```go
func (c *Command) DisableHelp() *Command
```
Disables automatic help flag for this command.

```go
func (c *Command) EnableHelp() *Command
```
Re-enables automatic help flag.

```go
func (c *Command) SetHelpFlag(name, short string) *Command
```
Customizes the help flag name.

```go
func (c *Command) IsHelpEnabled() bool
```
Returns whether help is enabled.

## Flag

### Flag Methods

```go
func (f *Flag) GetNames() []string
```
Returns all flag names (long and short).

```go
func (f *Flag) GetType() string
```
Returns flag type as string ("string", "int", "bool", etc.).

```go
func (f *Flag) PrimaryName() string
```
Returns the long flag name.

```go
func (f *Flag) ShortName() string
```
Returns the short flag name (empty if none).

```go
func (f *Flag) GetDefault() interface{}
```
Returns the default value.

```go
func (f *Flag) GetUsage() string
```
Returns the usage/help text.

```go
func (f *Flag) GetValue() interface{}
```
Returns the current flag value.

```go
func (f *Flag) IsRequired() bool
```
Returns whether flag is required.

```go
func (f *Flag) IsHidden() bool
```
Returns whether flag is hidden.

```go
func (f *Flag) IsSet() bool
```
Returns whether flag was explicitly set by user.

```go
func (f *Flag) HasName(name string) bool
```
Checks if flag has the given name (long or short).

## FlagSet

### Creating FlagSets

```go
func NewFlagSet() *FlagSet
```
Creates a new flag set (usually not needed, commands create their own).

### FlagSet Methods

```go
func (fs *FlagSet) Add(ptr interface{}, name, short string, defaultValue interface{}, usage string) *Flag
```
Adds a flag to the set.

```go
func (fs *FlagSet) GetFlag(name string) *Flag
```
Gets a flag by name (long or short).

```go
func (fs *FlagSet) GetFlags(names ...string) []*Flag
```
Gets multiple flags by name.

```go
func (fs *FlagSet) GetAll() []*Flag
```
Returns all flags in the set.

```go
func (fs *FlagSet) Parse(args []string) ([]string, error)
```
Parses flags from arguments, returns remaining args.

```go
func (fs *FlagSet) BindStruct(structPtr interface{})
```
Binds struct fields as flags using tags.

## Argument

### Argument Structure

```go
type Argument struct {
    Name        string  // Argument name
    Description string  // Help text
    Required    bool    // Whether required
}
```

## Completion

### Adding Completion

```go
func AddCompletion(rootCmd *Command)
```
Adds completion subcommands for all shells (bash, zsh, fish, powershell).

### Completion Types

```go
type BashCompletion struct{}
type ZshCompletion struct{}
type FishCompletion struct{}
type PowerShellCompletion struct{}
```

Each implements:
```go
func (c *Completion) GetCompletions(cmd *Command, args []string) []string
func (c *Completion) Register(cmd *Command)
func (c *Completion) GenerateScript(cmd *Command) string
```

## Errors

### Custom Error Types

```go
type CommandNotFoundError struct {
    Name string
}

type ArgumentError struct {
    Message string
}

type FlagError struct {
    Message string
}
```

All implement `error` interface with custom `Error()` messages.

## Type Conversions

### Supported Types

Flags and arguments support automatic type conversion for:

- **Strings**: `string`
- **Integers**: `int`, `int8`, `int16`, `int32`, `int64`
- **Unsigned**: `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- **Floats**: `float32`, `float64`
- **Boolean**: `bool`
- **Duration**: `time.Duration`
- **Arrays**: `[]string`, `[]int`, etc. (via repeated flags)

### Examples

```go
// String
var name string
cmd.Flag(&name, "name", "n", "default", "Name")

// Integer
var port int
cmd.Flag(&port, "port", "p", 8080, "Port")

// Boolean
var verbose bool
cmd.Flag(&verbose, "verbose", "v", false, "Verbose")

// Duration
var timeout time.Duration
cmd.Flag(&timeout, "timeout", "t", 30*time.Second, "Timeout")

// Array
var tags []string
cmd.Flag(&tags, "tag", "", nil, "Tags")
```

## Best Practices

### Error Handling

Always return errors from actions:

```go
cmd.Action(func(ctx context.Context, cmd *Command) error {
    if err := doWork(); err != nil {
        return fmt.Errorf("work failed: %w", err)
    }
    return nil
})
```

### Context Usage

Use context for cancellation:

```go
cmd.Action(func(ctx context.Context, cmd *Command) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case result := <-doWork():
        return nil
    }
})
```

### Flag Naming

- Use lowercase with hyphens: `--api-key`
- Provide short forms for common flags: `-v`, `-p`
- Be consistent across commands

### Command Naming

- Use lowercase: `deploy` not `Deploy`
- Use verbs for actions: `start`, `stop`, `deploy`
- Group related commands: `database migrate`, `database rollback`

## Examples

See the [examples/](../examples/) directory for complete working examples of all features.

## Next Steps

- Start with [Quick Start](quick-start.md)
- Learn about [Commands](commands.md)
- Explore [Flags and Arguments](flags-and-arguments.md)
- Understand [Lifecycle Hooks](lifecycle.md)
