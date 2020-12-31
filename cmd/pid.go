package cmd

import (
	"log"

	"github.com/mbreese/rtun/client"

	"github.com/spf13/cobra"
)

func init() {
	pidCmd.Flags().StringVarP(&socketFilename, "socket", "s", "", "Server socket")
	rootCmd.AddCommand(pidCmd)
}

var pidCmd = &cobra.Command{
	Use:    "pid",
	Short:  "Get the PID for the running server",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		client := client.Connect(socketFilename)

		defer client.Close()

		err := client.PID()
		if err != nil {
			log.Fatal(err)
		}

	},
}
