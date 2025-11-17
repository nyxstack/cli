# Flag Syntax

## Why `--flag=value` Instead of `--flag value`?

nyxstack/cli requires the `--flag=value` syntax for all non-boolean flags. This document explains why.

## The Problem with Space-Separated Syntax

Consider this command:

```bash
myapp --count 10 deploy
```

**Question**: What is `10`?
- Is it the value for `--count`?
- Or is it a positional argument before the `deploy` subcommand?

With space-separated syntax, it's ambiguous. The parser must make assumptions:

1. **Look-ahead parsing**: Check if `deploy` is a subcommand
2. **Type checking**: Check if `10` matches the flag type
3. **Argument counting**: Count remaining arguments

This creates complexity and edge cases.

## Real-World Example

```bash
# What does this mean?
myapp --env production deploy api

# Possibility 1: --env takes "production", deploy subcommand with "api" argument
# Possibility 2: --env takes "production deploy", "api" is the subcommand
# Possibility 3: --env takes no value (boolean), "production" is subcommand, "deploy" and "api" are args
```

## The Solution: `--flag=value`

With equals syntax, there's no ambiguity:

```bash
myapp --count=10 deploy
```

Now it's crystal clear:
- `--count` has value `10`
- `deploy` is a subcommand

## Syntax Rules

### Value Flags

Must use `=`:

```bash
# ✅ Correct
myapp --port=8080
myapp --name=john
myapp --count=10

# ❌ Incorrect
myapp --port 8080    # Error: flag needs a value
```

### Boolean Flags

Can omit value (defaults to `true`):

```bash
# ✅ All valid
myapp --verbose             # verbose = true
myapp --verbose=true        # verbose = true
myapp --verbose=false       # verbose = false
myapp -v                    # verbose = true (short form)
myapp -v=true               # verbose = true
```

### Short Flags

Same rule - must use `=` for value flags:

```bash
# ✅ Correct
myapp -p=8080
myapp -n=john
myapp -v              # Boolean, no value needed

# ❌ Incorrect
myapp -p 8080         # Error
```

## Benefits

### 1. No Ambiguity

```bash
# Clear separation between flags, subcommands, and arguments
myapp --env=prod --count=10 deploy api v2.1.0
      └─flag────┘ └─flag───┘ └─cmd─┘ └─args────┘
```

### 2. Simple Parsing

Parser can split on `=` and knows exactly what's what:

```go
parts := strings.Split("--port=8080", "=")
flag := parts[0]   // "--port"
value := parts[1]  // "8080"
```

### 3. Consistent Behavior

No special cases or look-ahead required:

```bash
# Works the same regardless of command structure
myapp --flag=value command
myapp command --flag=value
```

### 4. Better Error Messages

```bash
myapp --port
# Error: flag --port requires a value (use --port=VALUE)

myapp --port 8080
# Error: unknown flag or argument: --port
#        Did you mean --port=8080?
```

## Comparison with Other Frameworks

| Framework | Syntax | Ambiguity | Complexity |
|-----------|--------|-----------|------------|
| nyxstack/cli | `--flag=value` | None | Low |
| Cobra | `--flag value` | High | High |
| urfave/cli | `--flag value` | High | High |
| Docker | `--flag=value` | None | Low |
| kubectl | `--flag=value` | None | Low |

Note: Many modern CLI tools (Docker, kubectl, Kubernetes) use `--flag=value` syntax for the same reasons.

## Migration from Space-Separated

If you're migrating from a framework that uses `--flag value`:

```bash
# Before (Cobra/urfave)
myapp --port 8080 --host localhost deploy

# After (nyxstack/cli)
myapp --port=8080 --host=localhost deploy
```

Simple find-and-replace in documentation and examples.

## Edge Cases Eliminated

### Multiple Values

```bash
# Space-separated (ambiguous)
myapp --tag prod api v2    # What's a tag? What's an argument?

# Equals syntax (clear)
myapp --tag=prod --tag=api --tag=v2
```

### Flag After Subcommand

```bash
# Space-separated (complex parsing)
myapp deploy --env production    # Must parse "deploy" first

# Equals syntax (simple parsing)
myapp deploy --env=production    # Same complexity either way
```

### Boolean vs Value Flags

```bash
# Space-separated (must know flag type)
myapp --debug --port 8080    # Parser must know --debug is boolean

# Equals syntax (no type knowledge needed)
myapp --debug --port=8080    # Clear from syntax
```

## Real-World Usage

### Docker

```bash
docker run --name=myapp --port=8080:80 nginx
```

### Kubernetes

```bash
kubectl get pods --namespace=production --selector=app=api
```

### Git (mixed, but trending to =)

```bash
git commit --message="Fix bug"
git log --since="2 weeks ago"
```

## Best Practices

### Documentation

Always show the equals sign in examples:

```bash
# Good
myapp --port=8080

# Bad (ambiguous)
myapp --port 8080
```

### Help Text

```
Flags:
  --port=INT         Server port (required)
  --host=STRING      Server host (default: localhost)
  --verbose          Enable verbose output
```

### Error Messages

Be explicit about syntax:

```
Error: flag --port requires a value
Usage: --port=VALUE
```

## Frequently Asked Questions

### Q: Can I use spaces?

No. The syntax is intentionally strict to eliminate ambiguity.

### Q: Why not support both?

Supporting both syntax styles doubles the complexity and reintroduces ambiguity. It's better to have one clear way.

### Q: What about POSIX standards?

POSIX defines short flags with spaces (`-p 8080`), but doesn't define long flags. The `--flag=value` syntax is a GNU extension that's become standard in modern CLIs.

### Q: Is this a breaking change?

If you're migrating from space-separated syntax, yes. But it's a one-time change that eliminates an entire class of parsing bugs.

## Summary

The `--flag=value` syntax:
- ✅ Eliminates ambiguity
- ✅ Simplifies parsing
- ✅ Provides better error messages
- ✅ Matches modern CLI tools (Docker, kubectl)
- ✅ Makes complex commands clear

It's a small syntax change that makes CLI applications more reliable and easier to use.

## Next Steps

- See [Flags and Arguments](flags-and-arguments.md) for usage
- Read [Quick Start](quick-start.md) for examples
- Check [API Reference](api-reference.md) for details
