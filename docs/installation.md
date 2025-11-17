# Installation

## Requirements

- Go 1.20 or higher
- No external dependencies (except `github.com/nyxstack/color` for terminal colors)

## Install via go get

```bash
go get github.com/nyxstack/cli
```

## Verify Installation

Create a simple test program:

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/nyxstack/cli"
)

func main() {
    app := cli.Root("test").
        Action(func(ctx context.Context, cmd *cli.Command) error {
            fmt.Println("Installation successful!")
            return nil
        })
    
    app.Execute()
}
```

Run it:

```bash
go run main.go
# Installation successful!
```

## Using in Your Project

1. Initialize a new Go module (if you haven't already):

```bash
mkdir myapp
cd myapp
go mod init github.com/yourusername/myapp
```

2. Install nyxstack/cli:

```bash
go get github.com/nyxstack/cli
```

3. Create your CLI application in `main.go`

4. Build your binary:

```bash
go build -o myapp
./myapp --help
```

## Next Steps

- Follow the [Quick Start Guide](quick-start.md)
- Explore the [Examples](../examples/)
