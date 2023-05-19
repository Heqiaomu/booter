package cmd

import "github.com/spf13/cobra"

// CommandHooker for hook app's cmd impl
type CommandHooker interface {
	HookCmd(cmd *cobra.Command) error
}
