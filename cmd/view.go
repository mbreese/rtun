package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	viewCmd.Flags().StringVarP(&socketFilename, "socket", "s", "", "Socket filename (default $HOME/.rtun/rtun.sock.*)")
	rootCmd.AddCommand(viewCmd)
}

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View a file on the remote server",
	Run: func(cmd *cobra.Command, args []string) {
		var local string

		if len(args) > 0 {
			local = args[0]
		}

		finfo, err1 := os.Stat(local)
		if err1 != nil {
			log.Fatal(err1)
		}
		if finfo.IsDir() {
			log.Fatal("Cannot view a directory")
		}

		client := connect()
		defer client.Close()

		fmt.Printf("Viewing file: %s\n", local)

		err := client.View(local)
		if err != nil {
			fmt.Println(err1.Error())
			log.Fatal(err)
		}

	},
}
