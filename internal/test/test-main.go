package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// signal handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	select {
	case s := <-c:
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			return
		case syscall.SIGHUP:
		// TODO app reload
		default:
			return
		}
		return
	case <-time.After(30 * time.Second):
		return
	}
}
