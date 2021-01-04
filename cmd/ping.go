package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	pingCmd.Flags().StringVarP(&socketFilename, "socket", "s", "", "Socket filename (default $HOME/.rtun/rtun.sock.*)")
	rootCmd.AddCommand(pingCmd)
}

var pingCmd = &cobra.Command{
	Use:    "ping",
	Short:  "Ping the rtun server",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		client := connect()

		defer client.Close()

		ret, err := client.Ping()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(ret)

	},
}
