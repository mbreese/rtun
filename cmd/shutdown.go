package cmd

import (
	"errors"
	"log"

	"github.com/mbreese/rtun/client"

	"github.com/spf13/cobra"
)

func init() {
	shutdownCmd.Flags().StringVarP(&socketFilename, "socket", "s", "", "Server socket")
	shutdownCmd.MarkFlagRequired("socket")
	rootCmd.AddCommand(shutdownCmd)
}

var shutdownCmd = &cobra.Command{
	Use:    "shutdown",
	Short:  "Shutdown the rtun server",
	Hidden: true,

	Args: func(cmd *cobra.Command, args []string) error {
		if socketFilename == "" {
			return errors.New("Missing socket")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		client := client.Connect(socketFilename)

		defer client.Close()

		err := client.Shutdown()
		if err != nil {
			log.Fatal(err)
		}

	},
}
