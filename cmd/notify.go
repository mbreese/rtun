package cmd

import (
	"log"
	"strings"

	"github.com/mbreese/rtun/client"

	"github.com/spf13/cobra"
)

func init() {
	notifyCmd.Flags().StringVarP(&socketFilename, "socket", "s", "", "Server socket")
	notifyCmd.MarkFlagRequired("socket")
	rootCmd.AddCommand(notifyCmd)
}

var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "Send a notification to the server",
	Run: func(cmd *cobra.Command, args []string) {
		client := client.Connect(socketFilename)

		defer client.Close()

		err := client.Notify(strings.Join(args, " "))
		if err != nil {
			log.Fatal(err)
		}

	},
}
