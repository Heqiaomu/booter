package internal

import (
	"github.com/spf13/cobra"
)

// Start app
var Start = &cobra.Command{
	Use:   "start [App ...]",
	Short: "Start blocface app",
	Long: `
Start the blocface app and run until a stop command is received.
	`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return getProcess(cmd).start(cmd.Flag("nodaemon").Value.String() == "true", args)
	},
}

// Stop app
var Stop = &cobra.Command{
	Use:   "stop [App ...]",
	Short: "Stop blocface app",
	Long: `
Stop blocface app, if app is already stopped, return app's status.
`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return getProcess(cmd).stop(args)
	},
}

// Restart app
var Restart = &cobra.Command{
	Use:   "restart [App ...]",
	Short: "restart blocface app",
	Long: `
Restart the blocface app and run until a stop command is received.
	`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return getProcess(cmd).restart(args)
	},
}

// Status app
var Status = &cobra.Command{
	Use:   "status [App ...]",
	Short: "Dump status of app",
	Long: `
Dump status, check blocface app is running.
`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return getProcess(cmd).status(args)
	},
}

// Tail app log
var Tail = &cobra.Command{
	Use:   "tail [App ...]",
	Short: "Tail app log",
	Long: `
Display the log of the blocface app.
`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return getProcess(cmd).tail(args)
	},
}
