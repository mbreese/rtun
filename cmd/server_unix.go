// go:build linux freebsd darwin
// +build linux freebsd darwin !windows

package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"syscall"
)

// Fork crete a new process
// see: https://github.com/immortal/immortal/blob/master/fork.go
func fork(stdout string) (int, error) {
	args := make([]string, len(os.Args)-2)
	forkme := false
	for i, arg := range os.Args[1:] {
		if arg == "-d" || arg == "--daemon" {
			forkme = true
		} else if forkme {
			args[i-1] = arg
		} else {
			args[i] = arg
		}
	}

	// fmt.Printf("%v\n", args)

	var stdoutf *os.File
	if stdout == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		os.MkdirAll(path.Join(home, ".rtun"), 0700)
		stdout = path.Join(home, ".rtun", "rtun.log")
	}

	var err error
	stdoutf, err = os.OpenFile(stdout, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0700)
	if err != nil {
		fmt.Printf("Error setting up a file to write to stdout: %s\n", stdout)
		log.Fatal(err)
	}

	cmd := exec.Command(os.Args[0], args...)
	cmd.Env = os.Environ()
	cmd.Stdin = nil
	cmd.Stdout = stdoutf
	cmd.Stderr = stdoutf
	cmd.ExtraFiles = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	if err := cmd.Start(); err != nil {
		return 0, err
	}
	return cmd.Process.Pid, nil
}
