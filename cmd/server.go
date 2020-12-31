package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/mbreese/rtun/server"
)

var daemonize = false
var verbose = false
var stdout string
var downloadDir string

func init() {
	serverCmd.Flags().BoolVarP(&daemonize, "daemon", "d", false, "Run in the background")
	serverCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	serverCmd.Flags().StringVarP(&stdout, "stdout", "o", "", "Write output log to this file")
	serverCmd.Flags().StringVarP(&downloadDir, "dir", "", ".", "Save downloads to this directory")
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server <socket_file>",
	Short: "Starts the rtun server",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 || args[0] == "-" {
			// fmt.Printf("%v\n", args)
			return fmt.Errorf("missing local socket filename")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if daemonize {
			pid, _ := fork(stdout)
			fmt.Printf("%d\n", pid)
			os.Exit(0)
		}

		svr := server.NewServer(args[0], downloadDir, verbose)
		svr.Listen()
	},
}

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
	if stdout != "" {
		var err error
		stdoutf, err = os.OpenFile(stdout, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0700)
		if err != nil {
			fmt.Printf("Error setting up a file to write to stdout: %s\n", stdout)
			log.Fatal(err)
		}
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
