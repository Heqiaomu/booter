package internal

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Application base app framework interface
type Application interface {
	Init(ctx Context) error
	Run(ctx Context) error
	Exit(ctx Context) error

	// cmd func

	// other...
}

// Context is main application context
type Context struct {
	// ConfigReadMode is one of `local` or `etcd`
	ConfigReadMode string
	EtcdHost       string
	AppName        string
	KongHost       string
	// add others...

	// Magic context for hack
	Magic map[string]interface{}
}

// App is blocface main service
type App struct {
	Ctx *Context
}

// NewApp return a new blocface app
func NewApp() *App {
	return &App{}
}

// HookCmd will hook all commands default
// return if exit after cmd
func (app *App) HookCmd(cmd *cobra.Command) (err error) {

	cmd.PersistentFlags().String("work", "", "target app working dir")

	Start.PersistentFlags().Bool("nodaemon", false, "run in foreground")
	cmd.AddCommand(Start)

	cmd.AddCommand(Stop)
	cmd.AddCommand(Restart)
	cmd.AddCommand(Status)
	cmd.AddCommand(Tail)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		return err
	}

	return
}

// Init application init
func (app *App) Init() (err error) {

	return
}

// Run application main loop
func (app *App) Run() (err error) {

	return
}

// Exit application exit
func (app *App) Exit() (err error) {

	fmt.Println("booster app exit")

	return
}
