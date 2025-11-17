package cli

// ShellCompletion interface for different shell implementations
type ShellCompletion interface {
	// GetCompletions returns completion suggestions for a command
	GetCompletions(cmd *Command, args []string) []string

	// Register adds the completion command to the given command
	Register(cmd *Command)

	// GenerateScript generates the shell completion script
	GenerateScript(cmd *Command) string
}

// AddCompletion registers completion commands for all supported shells
func AddCompletion(rootCmd *Command) {
	// Register bash completion
	bashComp := &BashCompletion{}
	bashComp.Register(rootCmd)

	// Register zsh completion
	zshComp := &ZshCompletion{}
	zshComp.Register(rootCmd)

	// Register fish completion
	fishComp := &FishCompletion{}
	fishComp.Register(rootCmd)

	// Register PowerShell completion
	psComp := &PowerShellCompletion{}
	psComp.Register(rootCmd)
}

// getCompletionWords returns completion words for a command (shared implementation)
func getCompletionWords(cmd *Command) []string {
	var words []string

	// Add visible subcommands
	for name, subcmd := range cmd.GetCommands() {
		if !subcmd.IsHidden() {
			words = append(words, name)
		}
	}

	// Add all available flags
	allFlags := cmd.getAllFlags()
	for _, flag := range allFlags {
		// Skip hidden flags
		if flag.IsHidden() {
			continue
		}

		// Add primary name with --
		words = append(words, "--"+flag.PrimaryName())

		// Add short name with - if it exists
		if shortName := flag.ShortName(); shortName != "" {
			words = append(words, "-"+shortName)
		}
	}

	return words
}
