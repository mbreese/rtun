package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	echoCmd.Flags().StringVarP(&socketFilename, "socket", "s", "", "Socket filename (default $HOME/.rtun/rtun.sock.*)")
	rootCmd.AddCommand(echoCmd)
}

var echoCmd = &cobra.Command{
	Use:    "echo",
	Short:  "Echo a message to/from the server",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		client := connect()
		defer client.Close()

		ret, err := client.Echo(strings.Join(args, " "))
		if err != nil {
			log.Fatal(err)
		}

		if ret[:3] == "OK " {
			fmt.Println(ret[3:])
		} else {
			fmt.Println(ret)
		}

	},
}
