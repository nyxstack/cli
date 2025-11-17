package cli

// Argument represents a positional argument for a command
type Argument struct {
	Name        string
	Description string
	Required    bool
}
