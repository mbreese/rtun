package cmd

import (
	"log"

	"github.com/mbreese/rtun/client"

	"github.com/spf13/cobra"
)

func init() {
	pingCmd.Flags().StringVarP(&socketFilename, "socket", "s", "", "Server socket")
	pingCmd.MarkFlagRequired("socket")
	rootCmd.AddCommand(pingCmd)
}

var pingCmd = &cobra.Command{
	Use:    "ping",
	Short:  "Ping the rtun server",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		client := client.Connect(socketFilename)

		defer client.Close()

		err := client.Ping()
		if err != nil {
			log.Fatal(err)
		}

	},
}
