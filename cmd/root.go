package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/mbreese/rtun/client"
	"github.com/spf13/cobra"
)

var socketFilename string
var verbose = false

var (
	rootCmd = &cobra.Command{
		Use:     "rtun",
		Short:   "Reverse tunnel - a tunnel back to your local system (for SSH connections)",
		Version: "0.1.1",
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logging")
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func connect() *client.Client {

	if socketFilename == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		files, err := ioutil.ReadDir(path.Join(home, ".rtun"))

		for _, f := range files {
			if f.IsDir() {
				continue
			}

			if len(f.Name()) > 8 && f.Name()[:9] == "rtun.sock" {
				client, err := client.Connect(path.Join(home, ".rtun", f.Name()), verbose)

				if err != nil {
					os.Remove(path.Join(home, ".rtun", f.Name()))
					if verbose {
						fmt.Printf("Removing old socket file: %s\n", path.Join(home, ".rtun", f.Name()))
					}
					continue
				}

				return client

			}
		}
		panic(errors.New("Unable to find a valid server socket"))
	} else {
		client, err := client.Connect(socketFilename, verbose)
		if err != nil {
			panic(err)
		}
		return client
	}

}
