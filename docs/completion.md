# Shell Completion

nyxstack/cli provides built-in shell completion for Bash, Zsh, Fish, and PowerShell.

## Quick Start

Add completion to your root command:

```go
root := cli.Root("myapp").Description("My application")

// Add subcommands...
deploy := cli.Cmd("deploy")
rollback := cli.Cmd("rollback")
root.AddCommand(deploy)
root.AddCommand(rollback)

// Add completion support for all shells
cli.AddCompletion(root)
```

This automatically adds completion subcommands:
- `myapp completion bash`
- `myapp completion zsh`
- `myapp completion fish`
- `myapp completion powershell`

## Installation

### Bash

```bash
# Generate and source the completion script
myapp completion bash > /etc/bash_completion.d/myapp

# Or for current session only
source <(myapp completion bash)
```

Add to `~/.bashrc` for persistence:
```bash
source <(myapp completion bash)
```

### Zsh

```bash
# Generate completion script
myapp completion zsh > "${fpath[1]}/_myapp"

# Reload completions
compinit
```

Add to `~/.zshrc`:
```bash
autoload -U compinit; compinit
```

### Fish

```bash
myapp completion fish > ~/.config/fish/completions/myapp.fish
```

### PowerShell

```powershell
myapp completion powershell | Out-String | Invoke-Expression
```

Add to PowerShell profile:
```powershell
# Find profile location
$PROFILE

# Edit and add:
myapp completion powershell | Out-String | Invoke-Expression
```

## What Gets Completed

### Subcommands

```bash
myapp <TAB>
# Shows: deploy rollback help completion
```

### Flags

```bash
myapp --<TAB>
# Shows: --verbose --help

myapp deploy --<TAB>
# Shows: --env --verbose --help (includes inherited flags)
```

### Hidden Commands Excluded

Hidden commands don't appear in completion:

```go
debug := cli.Cmd("debug").Hidden()
root.AddCommand(debug)
```

```bash
myapp <TAB>
# Shows: deploy rollback help completion (debug is hidden)
```

## Completion Features

### Flag Name Completion

```bash
myapp --v<TAB>
# Completes to: myapp --verbose
```

### Short Flag Completion

```bash
myapp -<TAB>
# Shows: -v -h
```

### Subcommand Completion

```bash
myapp dep<TAB>
# Completes to: myapp deploy
```

### Multi-Level Completion

```bash
myapp database <TAB>
# Shows: migrate rollback backup

myapp database mig<TAB>
# Completes to: myapp database migrate
```

## Customization

### Manual Completion Registration

Instead of using `AddCompletion()`, register completions manually:

```go
bash := &cli.BashCompletion{}
bash.Register(root)

zsh := &cli.ZshCompletion{}
zsh.Register(root)

fish := &cli.FishCompletion{}
fish.Register(root)

ps := &cli.PowerShellCompletion{}
ps.Register(root)
```

### Custom Completion Command Name

```go
// Use "autocomplete" instead of "completion"
root.AddCommand(
    cli.Cmd("autocomplete").
        AddCommand(
            cli.Cmd("bash").Action(func(ctx context.Context, cmd *cli.Command) error {
                bash := &cli.BashCompletion{}
                fmt.Println(bash.GenerateScript(root))
                return nil
            }),
        ),
)
```

## Testing Completion

Test that completion works:

```bash
# Bash
complete -p myapp

# Zsh
echo $_comps[myapp]

# Fish
complete -C myapp

# PowerShell
Get-Command myapp -ShowCommandInfo
```

## How It Works

### Completion Flow

1. User types: `myapp dep<TAB>`
2. Shell calls: `myapp completion bash` (or zsh/fish/powershell)
3. Framework generates list of available commands/flags
4. Shell displays completions to user

### What's Included

- All visible subcommands (not hidden)
- All flags (including inherited)
- Short flag names
- Long flag names

### What's Excluded

- Hidden commands (`.Hidden()`)
- Hidden flags (`.FlagHidden()`)
- Arguments (positional parameters)

## Shell-Specific Features

### Bash

- Command completion
- Flag completion with `--` prefix
- Short flag completion with `-` prefix

### Zsh

- Advanced completion with descriptions
- Smart flag parsing
- Better handling of flag values

### Fish

- Real-time completion as you type
- Command descriptions
- Smart flag suggestions

### PowerShell

- Tab completion
- Command parameters
- Flag aliases

## Troubleshooting

### Completion Not Working

**Bash:**
```bash
# Check if completion is loaded
complete -p myapp

# Reload bash completion
source ~/.bashrc
```

**Zsh:**
```bash
# Check completion functions
echo $_comps[myapp]

# Rebuild completion cache
rm ~/.zcompdump*
compinit
```

**Fish:**
```bash
# Check completion file exists
ls ~/.config/fish/completions/myapp.fish

# Reload fish configuration
source ~/.config/fish/config.fish
```

**PowerShell:**
```powershell
# Check if completion is registered
Get-Command myapp -ShowCommandInfo

# Reload profile
. $PROFILE
```

### Completions Out of Date

Regenerate completion scripts after adding new commands:

```bash
# Bash
myapp completion bash > /etc/bash_completion.d/myapp

# Zsh
myapp completion zsh > "${fpath[1]}/_myapp"
compinit

# Fish
myapp completion fish > ~/.config/fish/completions/myapp.fish

# PowerShell (just reload)
. $PROFILE
```

## Best Practices

1. **Always call `AddCompletion()`**
   ```go
   cli.AddCompletion(root)  // Add this to every CLI app
   ```

2. **Document completion in help**
   ```go
   root.Description("My app\n\nShell completion: myapp completion bash|zsh|fish|powershell")
   ```

3. **Test completion during development**
   ```bash
   # Quick test
   eval "$(myapp completion bash)"
   myapp <TAB><TAB>
   ```

4. **Include setup instructions in README**
   ```markdown
   ## Shell Completion

   ```bash
   # Bash
   source <(myapp completion bash)

   # Zsh
   myapp completion zsh > "${fpath[1]}/_myapp" && compinit
   ```
   ```

## Examples

See [examples/completion/](../examples/completion/) for a complete working example.

## Comparison with Other Frameworks

| Feature | nyxstack/cli | Cobra | urfave/cli |
|---------|--------------|-------|------------|
| Bash | ✅ | ✅ | ⚠️ Basic |
| Zsh | ✅ | ✅ | ⚠️ Basic |
| Fish | ✅ | ✅ | ❌ |
| PowerShell | ✅ | ✅ | ❌ |
| Auto-registration | ✅ One call | ⚠️ Manual | ⚠️ Manual |
| Hidden commands | ✅ Excluded | ✅ Excluded | N/A |
| Flag inheritance | ✅ Automatic | ⚠️ Manual | ❌ |

## Next Steps

- Check the [Quick Start](quick-start.md) guide
- Learn about [Commands](commands.md)
- Explore [Flags and Arguments](flags-and-arguments.md)
