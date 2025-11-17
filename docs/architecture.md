# Architecture

Design principles and system architecture of nyxstack/cli.

## Design Principles

### 1. Simple by Default
Most common use cases should require minimal code:

```go
cli.Root("myapp").
    Flag(&verbose, "verbose", "v", false, "Verbose").
    Action(runApp).
    Execute()
```

### 2. Powerful When Needed
Advanced features available without adding complexity to simple cases:

```go
cmd.
    PersistentPreRun(initDB).
    PreRun(validate).
    Action(deploy).
    PostRun(cleanup).
    PersistentPostRun(closeDB)
```

### 3. Type Safety
Leverage Go's type system for correctness:

```go
Action(func(ctx context.Context, cmd *Command, name string, count int) error {
    // name is string, count is int - compiler enforces types
})
```

### 4. Zero Magic
Explicit is better than implicit:

```go
// Explicit flag definition
cmd.Flag(&port, "port", "p", 8080, "Server port")

// Explicit inheritance (automatic, but visible)
root.AddCommand(child)  // child inherits root's flags
```

### 5. Testable
Easy to test without running actual binaries:

```go
err := cmd.ExecuteWithArgs([]string{"--flag=value", "arg"})
```

### 6. Silent Framework
You control all output:

```go
// Framework never prints anything except help
// You decide what to print and when
fmt.Println("Deploying...")
```

## System Architecture

### Component Overview

```
┌─────────────────────────────────────────────┐
│              Application                     │
└─────────────────────────────────────────────┘
                     │
                     ├── Root Command
                     │   ├── Flags (inherited by all)
                     │   ├── Lifecycle Hooks
                     │   └── Subcommands
                     │       ├── Command A
                     │       │   ├── Flags (own + inherited)
                     │       │   ├── Arguments
                     │       │   └── Action
                     │       └── Command B
                     │           ├── Flags (own + inherited)
                     │           └── Subcommands
                     │               └── Command C
                     │
                     ├── Flag System
                     │   ├── FlagSet (per command)
                     │   ├── Flag Parser
                     │   └── Type Conversion
                     │
                     ├── Execution Engine
                     │   ├── Command Router
                     │   ├── Lifecycle Manager
                     │   └── Context Handler
                     │
                     └── Help System
                         ├── Auto-generation
                         └── Completion Scripts
```

### Core Components

#### 1. Command
The central abstraction representing a CLI command.

**Responsibilities:**
- Store command metadata (name, description)
- Manage flags and arguments
- Define lifecycle hooks
- Maintain parent-child relationships

**Key Methods:**
- `Flag()` - Add flags
- `Arg()` - Add arguments
- `Action()` - Define behavior
- `AddCommand()` - Build hierarchy

#### 2. FlagSet
Manages a collection of flags for a command.

**Responsibilities:**
- Store flags
- Parse flag values from arguments
- Inherit parent flags
- Validate required flags

**Features:**
- Pointer-based storage (`[]*Flag`)
- Set tracking (knows if flag was explicitly set)
- Type conversion

#### 3. Execution Engine
Orchestrates command execution.

**Flow:**
```
1. Parse arguments → Identify command
2. Parse flags → Set values
3. Validate → Check required flags/args
4. Run PersistentPreRun hooks (parent → child)
5. Run PreRun hook (this command)
6. Run Action (this command)
7. Run PostRun hook (this command)
8. Run PersistentPostRun hooks (child → parent)
```

#### 4. Type System
Handles automatic type conversion.

**Supported Types:**
- Primitives: string, int, bool, float64
- Extended: int8-64, uint8-64, float32, time.Duration
- Collections: []string, []int, etc.

**Conversion Points:**
- Flag values: `--count=10` → int
- Arguments: `myapp deploy 5` → int
- Action parameters: Automatic injection

## Data Flow

### Flag Parsing

```
Input: ["--port=8080", "deploy", "--env=prod"]
                 ↓
1. Split into flags and args
   Flags: ["--port=8080", "--env=prod"]
   Args:  ["deploy"]
                 ↓
2. Parse each flag
   --port=8080 → name="port", value="8080"
   --env=prod  → name="env", value="prod"
                 ↓
3. Find command
   "deploy" → Execute deploy command
                 ↓
4. Set flag values
   port = 8080 (int)
   env = "prod" (string)
```

### Flag Inheritance

```
Root: --verbose, --config
  ↓
Child: --env (inherits --verbose, --config)
  ↓
Grandchild: --port (inherits all three)

getAllFlags() → merges parent flags with own flags
```

### Action Execution

```
User: myapp deploy production 3
              ↓
1. Parse: command="deploy", args=["production", "3"]
              ↓
2. Match action signature:
   func(ctx, cmd, env string, replicas int) error
              ↓
3. Convert arguments:
   "production" → string (no conversion)
   "3" → int (convert)
              ↓
4. Call action with converted args
```

## Key Design Decisions

### 1. `--flag=value` Syntax

**Why:** Eliminates ambiguity between flag values and arguments.

**Alternative:** `--flag value` (rejected due to parsing complexity)

See [Flag Syntax](flag-syntax.md) for details.

### 2. Automatic Flag Inheritance

**Why:** Common flags (verbose, config) shouldn't be repeated on every command.

**Implementation:** Each command merges parent flags with its own.

### 3. Pointer-Based FlagSet

**Why:** Need to track if flag was explicitly set by user.

```go
type Flag struct {
    value interface{}
    set   bool  // Was flag provided by user?
}
```

### 4. Type-Safe Arguments

**Why:** Reduce runtime errors, leverage compile-time checking.

**Implementation:** Reflect on action function signature, convert arguments to match.

### 5. Context-First API

**Why:** Modern Go idiom for cancellation and timeouts.

```go
Action(func(ctx context.Context, cmd *Command) error)
```

### 6. Fluent API

**Why:** Readable, chainable, self-documenting.

```go
cli.Root("app").
    Flag(...).
    Arg(...).
    Action(...)
```

## Performance Considerations

### Flag Lookup
- O(n) linear search through flag list
- Acceptable for typical CLI (< 50 flags)
- Could optimize with map if needed

### Command Routing
- O(d) where d is depth of command hierarchy
- Typical depth: 1-3 levels
- No performance concerns

### Memory Usage
- One Command struct per command
- One Flag struct per flag
- Minimal overhead

## Extensibility

### Custom Completion
Implement `ShellCompletion` interface:

```go
type CustomCompletion struct{}

func (c *CustomCompletion) GetCompletions(cmd *Command, args []string) []string
func (c *CustomCompletion) Register(cmd *Command)
func (c *CustomCompletion) GenerateScript(cmd *Command) string
```

### Custom Types
Flags/arguments use `interface{}` with type switching:

```go
case MyCustomType:
    // Handle custom type
```

### Custom Validation
Use PreRun hooks:

```go
cmd.PreRun(func(ctx context.Context, cmd *Command) error {
    return validateMyRules()
})
```

## Future Considerations

### Potential Enhancements
1. **Custom help templates** - User-defined help formatting
2. **Command aliases** - Multiple names for same command
3. **Flag groups** - Mutually exclusive or required together
4. **Environment variable binding** - Auto-read from env vars
5. **Config file integration** - Load flags from config files

### Maintained Simplicity
Any enhancement must:
- Not complicate simple use cases
- Be opt-in, not required
- Have clear, documented API
- Maintain backward compatibility

## Comparison with Other Frameworks

### vs Cobra
- **Simpler API**: Fluent vs struct-based
- **Automatic inheritance**: Built-in vs manual
- **Type-safe args**: Yes vs no
- **Trade-off**: Less flexible, more opinionated

### vs urfave/cli
- **Better structure**: Hierarchical vs flat
- **Flag inheritance**: Automatic vs manual
- **Context support**: Built-in vs manual
- **Trade-off**: More complex implementation

## Next Steps

- Read [Commands](commands.md) for usage patterns
- See [Flags and Arguments](flags-and-arguments.md) for details
- Check [API Reference](api-reference.md) for complete API
