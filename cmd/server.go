package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/mbreese/rtun/client"
	"github.com/mbreese/rtun/server"
)

var daemonize = false
var stdout string
var downloadDir string

func init() {
	serverCmd.Flags().StringVarP(&socketFilename, "socket", "s", "", "Socket filename (default $HOME/.rtun/rtun.sock)")
	serverCmd.Flags().BoolVarP(&daemonize, "daemon", "d", false, "Run in the background")
	serverCmd.Flags().StringVarP(&stdout, "log", "l", "", "Write output log to this file (in daemon mode, default $HOME/.rtun/rtun.log)")
	serverCmd.Flags().StringVarP(&downloadDir, "dir", "D", ".", "Save downloads to this directory")
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the rtun server",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if socketFilename == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				panic(err)
			}
			os.MkdirAll(path.Join(home, ".rtun"), 0700)
			socketFilename = path.Join(home, ".rtun", "rtun.sock")
		}

		_, err := os.Stat(socketFilename)
		if err == nil || !os.IsNotExist(err) {

			// check to see if the server is still active. If not, remove the socketFile and keep going.
			_, err := client.Connect(socketFilename, verbose)
			if err != nil {
				// bad client.
				fmt.Fprintf(os.Stderr, "Socket file already in use: %s, but server not running. Starting new server.\n", socketFilename)
				os.Remove(socketFilename)
			} else {
				fmt.Fprintf(os.Stderr, "Socket file already in use: %s\n", socketFilename)
				os.Exit(1)
			}
		}

		if daemonize {
			pid, err := fork(stdout)
			if err == nil {
				fmt.Printf("%d\n", pid)
				os.Exit(0)
			} else {
				fmt.Printf("%s\n", err.Error())
				os.Exit(1)
			}
		}

		svr := server.NewServer(socketFilename, downloadDir, verbose)
		svr.Listen()

	},
}
