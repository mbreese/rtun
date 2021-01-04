package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	notifyCmd.Flags().StringVarP(&socketFilename, "socket", "s", "", "Socket filename (default $HOME/.rtun/rtun.sock.*)")
	rootCmd.AddCommand(notifyCmd)
}

var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "Send a notification to the server",
	Run: func(cmd *cobra.Command, args []string) {
		client := connect()

		defer client.Close()

		ret, err := client.Notify(strings.Join(args, " "))
		if err != nil {
			log.Fatal(err)
		}

		if ret != "OK" {
			log.Fatal("Unknown error")
		}

	},
}
