package main

import (
	"context"
	"fmt"

	"github.com/nyxstack/cli"
)

var rootCmd = cli.Cmd("myapp").
	Description("My application").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		cmd.ShowHelp()
		return nil
	})

var deployCmd = cli.Cmd("deploy").
	Description("Deploy the application").
	Arg("environment", "Target environment", true).
	Action(func(ctx context.Context, cmd *cli.Command, env string) error {
		fmt.Printf("Deploying to %s\n", env)
		return nil
	})
var serverCmd = cli.Cmd("server").
	Description("Server management").
	Action(func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("Server command")
		return nil
	})

func init() {
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(serverCmd)
}
func main() {
	fmt.Println("starting cli.")
	fmt.Println(rootCmd.Execute())
	// This is a placeholder main function.
	// The actual implementation is in the other files.
}
