package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	echoCmd.Flags().StringVarP(&socketFilename, "socket", "s", "", "Server socket")
	echoCmd.MarkFlagRequired("socket")
	rootCmd.AddCommand(echoCmd)
}

var echoCmd = &cobra.Command{
	Use:    "echo",
	Short:  "Echo a message to/from the server",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		client := connect()
		defer client.Close()

		err := client.Echo(strings.Join(args, " "))
		if err != nil {
			log.Fatal(err)
		}

	},
}
