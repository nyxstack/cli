# CLI Framework Examples

This directory contains example programs demonstrating various features of the CLI framework.

## Examples

### basic
Basic CLI with root command, flags, and a simple action.

```bash
cd basic
go run main.go --help
go run main.go --verbose --config=app.yaml
```

### subcommands
Application with nested subcommands and argument handling.

```bash
cd subcommands
go run main.go --help
go run main.go deploy staging v1.2.3 --verbose
go run main.go server start --port=3000 --detach
go run main.go server stop
```

### flags-struct
Demonstrates struct-based flag definition using tags.

```bash
cd flags-struct
go run main.go --help
go run main.go --host=0.0.0.0 --port=3000 --tls --timeout=1m
```

### lifecycle
Shows lifecycle hooks execution order (PersistentPreRun, PreRun, Action, PostRun, PersistentPostRun).

```bash
cd lifecycle
go run main.go --verbose
go run main.go sub --verbose
```

### completion
Shell completion support with hidden commands.

```bash
cd completion
go run main.go completion bash > /tmp/completion.sh
source /tmp/completion.sh

# Try tab completion
go run main.go <TAB>
go run main.go database <TAB>
```

### validation
Input validation using PreRun hooks.

```bash
cd validation
go run main.go --format=xml  # Error: invalid format
go run main.go --port=70000  # Error: port out of range
go run main.go --format=json --port=8080  # Success
```

### inheritance
Demonstrates automatic flag inheritance from parent to child commands.

```bash
cd inheritance
go run main.go --verbose --config=app.yaml
go run main.go server --verbose --host=0.0.0.0
go run main.go server start --verbose --port=3000 --detach
```

### hidden
Shows hidden flags and hidden commands that don't appear in help/completions.

```bash
cd hidden
go run main.go --help          # Won't show --debug or --api-key flags
go run main.go                 # Won't show 'secret' command
go run main.go public          # Public command works
go run main.go secret          # Secret command still executable
go run main.go --debug --api-key=test123  # Hidden flags still work
```

### argument-types
Demonstrates automatic type conversion for arguments (string, int, bool, float64, time.Duration).

```bash
cd argument-types
go run main.go greet Alice
go run main.go repeat "Hello" 3
go run main.go toggle debugging true
go run main.go calculate 3.14 2.5
go run main.go wait 2s
go run main.go config 8080 30s 5
```

### context-usage
Shows context usage patterns: timeouts, cancellation, and storing request-scoped values.

```bash
cd context-usage
go run main.go timeout --timeout=3s
go run main.go cancel
go run main.go values --user=john --log-level=debug
go run main.go combined --timeout=5s --user=alice
```

### array-flags
Demonstrates array/slice flags that accept multiple values.

```bash
cd array-flags
go run main.go deploy production --tag=v1.2.3 --tag=stable --tag=prod
go run main.go search "*.go" --include="src/**" --include="lib/**" --exclude="test/**"
go run main.go server --host=localhost --host=0.0.0.0 --host=192.168.1.1
go run main.go run myapp --env=PORT=8080 --env=DEBUG=true --env=API_KEY=secret
```

### flag-shadowing
Shows how child commands can override (shadow) parent command flags.

```bash
cd flag-shadowing
go run main.go --timeout=10 query "SELECT *"        # Uses root timeout: 10
go run main.go database --timeout=120               # Uses db timeout: 120
go run main.go server --verbose --format=json       # Server shadows verbose/format
go run main.go database start --timeout=90          # Uses db shadowed timeout
```

### error-handling
Demonstrates error handling in lifecycle hooks and execution flow control.

```bash
cd error-handling
go run main.go validate test123                     # Success
go run main.go validate test --fail-at=prerun       # PreRun error stops execution
go run main.go process --fail-at=action             # Action error, but PostRun still runs
go run main.go cleanup --fail-at=postrun            # PostRun error after successful action
go run main.go chain --fail-at=persistent-prerun    # Early error stops chain
go run main.go success                              # All hooks execute successfully
```

### execute-context
Shows Execute vs ExecuteContext patterns for different use cases.

```bash
cd execute-context
go run main.go basic                                # Basic Execute pattern
go run main.go fail                                 # Error handling
go run main.go long --timeout=3s                    # Context timeout (will cancel)
go run main.go long --timeout=15s                   # Context timeout (will complete)
go run main.go server                               # Graceful shutdown with cancellation
```

## Running Examples

Each example is a standalone Go program. Navigate to the example directory and run:

```bash
go run main.go [arguments...]
```

Add `--help` to any command to see usage information.

## Building Examples

To build an example:

```bash
cd <example-name>
go build -o app
./app --help
```
