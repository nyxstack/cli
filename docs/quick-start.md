# Quick Start Guide

Get started with nyxstack/cli in 5 minutes.

## Installation

```bash
go get github.com/nyxstack/cli
```

## Your First CLI App

Create a file called `main.go`:

```go
package main

import (
    "context"
    "fmt"
    "os"
    
    "github.com/nyxstack/cli"
)

func main() {
    var verbose bool
    
    app := cli.Root("myapp").
        Description("My awesome CLI application").
        Flag(&verbose, "verbose", "v", false, "Enable verbose output").
        Action(func(ctx context.Context, cmd *cli.Command) error {
            if verbose {
                fmt.Println("Verbose mode enabled")
            }
            fmt.Println("Hello from myapp!")
            return nil
        })
    
    if err := app.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

Run it:

```bash
go run main.go
# Hello from myapp!

go run main.go --verbose=true
# Verbose mode enabled
# Hello from myapp!

go run main.go -v
# Verbose mode enabled
# Hello from myapp!

go run main.go --help
# Shows automatic help
```

## Adding Subcommands

Add a `deploy` subcommand:

```go
func main() {
    var verbose bool
    
    root := cli.Root("myapp").
        Description("My awesome CLI application").
        Flag(&verbose, "verbose", "v", false, "Enable verbose output")
    
    // Add deploy subcommand
    var env string
    deploy := cli.Cmd("deploy").
        Description("Deploy the application").
        Flag(&env, "env", "e", "production", "Target environment").
        Action(func(ctx context.Context, cmd *cli.Command) error {
            if verbose {
                fmt.Printf("Deploying to %s (verbose mode)\n", env)
            } else {
                fmt.Printf("Deploying to %s\n", env)
            }
            return nil
        })
    
    root.AddCommand(deploy)
    
    if err := root.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

Run it:

```bash
myapp deploy
# Deploying to production

myapp deploy --env=staging
# Deploying to staging

myapp deploy --verbose=true --env=dev
# Deploying to dev (verbose mode)
# Note: --verbose flag inherited from parent!
```

## Adding Arguments

Add positional arguments to your command:

```go
deploy := cli.Cmd("deploy").
    Description("Deploy the application").
    Arg("service", "Service name to deploy", true).
    Arg("version", "Version to deploy", false).
    Action(func(ctx context.Context, cmd *cli.Command, service, version string) error {
        if version == "" {
            version = "latest"
        }
        fmt.Printf("Deploying %s:%s\n", service, version)
        return nil
    })
```

Run it:

```bash
myapp deploy api
# Deploying api:latest

myapp deploy api v2.1.0
# Deploying api:v2.1.0

myapp deploy
# Error: missing required argument: service
```

## Next Steps

- Learn about [Commands](commands.md)
- Explore [Flags and Arguments](flags-and-arguments.md)
- Understand [Lifecycle Hooks](lifecycle.md)
- Check out the [Examples](../examples/)
