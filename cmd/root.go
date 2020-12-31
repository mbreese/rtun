package cmd

import (
	"github.com/mbreese/rtun/client"
	"github.com/spf13/cobra"
)

var socketFilename string

var (
	rootCmd = &cobra.Command{
		Use:     "rtun",
		Short:   "Reverse tunnel - a tunnel back to your local system (for SSH connections)",
		Version: "0.1.0",
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func connect() *client.Client {
	client := client.Connect(socketFilename)
	return client
}
