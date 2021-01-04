package cmd

import (
	"errors"
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
		if socketFilename == "" {
			return errors.New("Missing socket")
		}

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
