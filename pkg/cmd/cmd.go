package cmd

import (
	"fmt"
	"github.com/Heqiaomu/booter/internal"
	"github.com/Heqiaomu/booter/pkg/run"
	"os"
	"path/filepath"
	"strings"

	"github.com/Heqiaomu/booter/pkg/version"
	util "github.com/Heqiaomu/goutil"

	"github.com/spf13/cobra"
)

var (
	appNameSplit = strings.Split(filepath.Base(os.Args[0]), "_")
	appName      = appNameSplit[len(appNameSplit)-1]
	exitAfterCmd = true
)

var (
	appPtr   *internal.Application // app Pointer
	vflag    bool                  // flag to show version info
	nodaemon bool                  // flag to tern on nodaemon mode
	confMode string                // flag to switch config read mode
	etcdHost string                // flag to read etcd host
	kongHost string                // flag to read kong host
	workdir  string                // working dir of bin
)

// PreBuild build default command for app
func PreBuild(app interface{}, ctx *internal.Context) (err error) {

	// base sub commands
	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().BoolVarP(&vflag, "version", "v", false, "version info output")
	rootCmd.PersistentFlags().StringVarP(&workdir, "dir", "", ".", "binary working dir")

	if hook, ok := app.(CommandHooker); ok {
		// hook cmd
		return hook.HookCmd(rootCmd)
	}

	// extra flags
	rootCmd.PersistentFlags().StringVarP(&confMode, "config", "c", "local", "'local' mode or 'etcd' mode")
	rootCmd.PersistentFlags().StringVarP(&kongHost, "kong", "", "http://127.0.0.1:8001", "kong host")
	rootCmd.PersistentFlags().StringVarP(&etcdHost, "etcd", "", "127.0.0.1:2379", "etcd host")

	// extra sub flags
	startCmd.PersistentFlags().BoolVarP(&nodaemon, "nodaemon", "", false, "open daemon mode")

	// add commands
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(statusCmd)

	// exec cmd
	return execute(app, ctx)
}

func execute(app interface{}, ctx *internal.Context) (err error) {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return err
	}
	// build context
	ctx.ConfigReadMode = confMode
	ctx.EtcdHost = etcdHost
	ctx.AppName = appName
	ctx.KongHost = kongHost

	// check arg
	if !util.StringInSlice(confMode, []string{"local", "etcd"}) {
		fmt.Println("config mode should be `local` or `etcd`")
		os.Exit(-1)
	}

	// handle --etcd flags
	if confMode == "etcd" && etcdHost == "" {
		fmt.Println("etcd host should not be empty")
		os.Exit(-1)
	}

	// handle --kong flags
	if kongHost == "" {
		fmt.Println("kong host should not be empty")
		os.Exit(-1)
	}

	// set work dir
	if workdir != "." {
		if e, _ := util.PathExists(workdir); !e {
			fmt.Printf("%s does not exists\n", workdir)
			os.Exit(-1)
		}
		absp, err := filepath.Abs(workdir)
		if err != nil {
			fmt.Println("set dir failed: ", err)
			os.Exit(-1)
		}
		fmt.Println("set workdir : ", absp)
		os.Chdir(absp)
	}

	// exit application when cmd mode
	if exitAfterCmd {
		os.Exit(0)
	}

	return nil
}

// root command
var rootCmd = &cobra.Command{
	Use:   appName,
	Short: appName + " is a blocface System sub service",
	Long: `
====================================================
Part of the command line is an auxiliary tool, which
is mainly used to start, stop, restart applications 
and child processes, and query application status.
====================================================`,
	Run: func(cmd *cobra.Command, args []string) {

		// handle -v --version flags
		if vflag {
			fmt.Println(vstr())
			os.Exit(0)
		}

		cmd.Help()
	},
}

// /////// commands //////////
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version info",
	Long:  "Print version information of " + appName,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(vstr())
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start application",
	Long:  "Start application and all sub processes",
	RunE: func(cmd *cobra.Command, args []string) error {
		// run no daemon
		if nodaemon {
			exitAfterCmd = false
		} else {
			// 后台运行app
			return run.RunAppWithName(appName)
		}

		return nil
	},
}

// restart when daemon mode
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart application",
	Long:  "Restart application and all sub processes",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if err = run.StopAppWithName(appName); err != nil {
			return err
		}
		if err = run.RunAppWithName(appName); err != nil {
			return err
		}

		return nil
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop application",
	Long:  "Stop application and all sub processes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run.StopAppWithName(appName)
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Query application status",
	Long:  "Query application status and all sub processes' status",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run.AppStatus(appName)
	},
}

var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Debug grpc",
	Long:  "Debug all grpc services",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO command: grpc")
	},
}

func vstr() (str string) {
	value, err := version.GetString()
	if err != nil {
		fmt.Println(err)
		return
	}
	str = value
	return
}
