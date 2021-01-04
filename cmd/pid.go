package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	pidCmd.Flags().StringVarP(&socketFilename, "socket", "s", "", "Socket filename (default $HOME/.rtun/rtun.sock.*)")
	rootCmd.AddCommand(pidCmd)
}

var pidCmd = &cobra.Command{
	Use:    "pid",
	Short:  "Get the PID for the running server",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		client := connect()

		defer client.Close()

		ret, err := client.PID()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(ret)

	},
}
