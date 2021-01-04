package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	shutdownCmd.Flags().StringVarP(&socketFilename, "socket", "s", "", "Socket filename (default $HOME/.rtun/rtun.sock.*)")
	rootCmd.AddCommand(shutdownCmd)
}

var shutdownCmd = &cobra.Command{
	Use:    "shutdown",
	Short:  "Shutdown the rtun server",
	Hidden: true,

	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		client := connect()

		defer client.Close()

		ret, err := client.Shutdown()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(ret)

	},
}
