package run

// ========================
// 不适用booter管理的流程
// 可以与booter一起使用
// ========================

import (
	"fmt"
	"github.com/Heqiaomu/booter/internal"
	"os"
	"os/exec"
	"path/filepath"
)

// RunAppWithName 通过app名运行app, 不经过booter统一调度
func RunAppWithName(bin string) (err error) {
	if r, pid := internal.Running(bin); r {
		fmt.Print("[", bin, "] ", pid, "\n")
		return nil
	}

	// app bin path
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	cmd := exec.Command(filepath.Join(dir, bin), []string{"start", "--nodaemon"}...)

	// no log file default
	if err = cmd.Start(); err != nil {
		return err
	}
	if s, pid := internal.Started(bin); s {
		fmt.Print("[", bin, "] ", pid, "\n")
	}

	return
}

// StopAppWithName 通过app名停止app, 不经过booter统一调度
func StopAppWithName(bin string) (err error) {
	run, pid := internal.Running(bin)
	if !run {
		fmt.Print("[", bin, "] DOWN\n")
		return nil
	}

	cmd := exec.Command("kill", "-TERM", pid)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err == nil {
		fmt.Print("[", bin, "] DOWN\n")
	}
	return
}

// AppStatus show app status
func AppStatus(bin string) (err error) {
	if r, pid := internal.Running(bin); r {

		fmt.Print("[", bin, "] ", pid, ", Started At: ", internal.StartAt(pid), "\n")
		return nil
	}
	fmt.Print("[", bin, "] DOWN\n")
	return nil
}
