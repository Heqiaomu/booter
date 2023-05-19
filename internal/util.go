package internal

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// https://github.com/shirou/gopsutil
// + https://github.com/foreversmart/sir
// + grep -f 实现独立的进程管理和监控模块,其他模块引用即可

// 判断是否正在运行
func Running(bin string) (run bool, pid string) {
	output, _ := exec.Command("pgrep", "-f", bin).Output()
	pidStr := strings.TrimSpace(string(output))

	return !(pidStr == ""), pidStr
}

// 启动app
func runapp(console bool, a app) (err error) {
	cmd := exec.Command(a.bin, a.args...)

	if console {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// 增加输出日志, TODO 日志分割
	logOutput, err := open(a.log)
	if err != nil {
		return err
	}
	defer logOutput.Close()
	cmd.Stdout = logOutput
	cmd.Stderr = logOutput
	return cmd.Start()
}

// 打开文件
func open(name string) (file *os.File, err error) {
	basep, _ := filepath.Split(name)
	if err := os.MkdirAll(basep, 0755); err != nil {
		return nil, err
	}

	logOutput, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return logOutput, nil
}

// 判断是否已经启动成功
func Started(bin string) (s bool, pid string) {
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if r, pid := Running(bin); r {
				return true, pid
			}
		case <-time.After(time.Second):
			return false, ""
		}
	}
}

// 获取程序启动时间
func StartAt(pid string) (start string) {
	shell := `ps -p ` + pid + ` -o lstart | awk 'NR==2{print $2,$3,$4,$5}'`
	output, _ := exec.Command("/bin/sh", "-c", shell).Output()
	start = strings.TrimSpace(string(output))

	return start
}
