package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/mbreese/rtun/client"

	"github.com/spf13/cobra"
)

var recurse bool

func init() {
	sendCmd.Flags().StringVarP(&socketFilename, "socket", "s", "", "Socket filename (default $HOME/.rtun/rtun.sock.*)")
	sendCmd.Flags().BoolVarP(&recurse, "recurse", "r", false, "Recursively upload directories")
	rootCmd.AddCommand(sendCmd)
}

var sendCmd = &cobra.Command{
	Use:   "send <local_file> [<remote_filename>]",
	Short: "Upload a file back to the local machine",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Missing file(s) to send...")
			cmd.Usage()
			os.Exit(1)
		}

		// var local string
		var remote string

		// if len(args) > 0 {
		// 	local = args[0]
		// }
		if len(args) > 1 {
			remote = args[len(args)-1]
		}
		client := connect()
		defer client.Close()

		for i, local := range args {
			if i > 0 && i == len(args)-1 {
				continue
			}
			fmt.Printf("Local: %s, Remote: %s\n", local, remote)

			finfo, err1 := os.Stat(local)
			if err1 != nil {
				fmt.Fprintf(os.Stderr, "File not found: %s\n", local)
				cmd.Usage()
				os.Exit(1)

				// fmt.Println(err1.Error())
				// log.Fatal(err1)
			}

			// fmt.Printf("Local: %s, Remote: %s\n", local, remote)

			if finfo.IsDir() {
				if !recurse {
					fmt.Fprintf(os.Stderr, "Not sending directory without --recurse/-r flag\n")
					cmd.Usage()
					os.Exit(1)
				}
				err := sendDir(local, remote, client)
				if err != nil {
					log.Fatal(err)
				}

			} else {
				remote2 := remote

				if remote != "" && remote[len(remote)-1] == os.PathSeparator {
					remote2 = path.Join(remote, path.Base(local))
				}

				err := client.Send(local, remote2)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	},
}

func sendDir(localDir string, remote string, client *client.Client) error {
	return innerSendDir(localDir, remote, "", client)
}

func innerSendDir(localDir string, remote string, curPath string, client *client.Client) error {
	files, err := ioutil.ReadDir(localDir)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	for _, f := range files {
		if f.IsDir() {
			// fmt.Printf("Recurse: %s\n", localDir+string(os.PathSeparator)+f.Name())
			innerSendDir(path.Join(localDir, f.Name()), remote, path.Join(curPath, f.Name()), client)
			continue
		}

		if !f.Mode().IsRegular() {
			fmt.Printf("Not sending (not a file): %s\n", path.Join(localDir, f.Name()))
			continue
		}

		remote1 := path.Clean(path.Join(remote, curPath, f.Name()))

		// fmt.Printf("SEND %s => %s\n", localDir+string(os.PathSeparator)+f.Name(), remote1)

		err := client.Send(path.Join(localDir, f.Name()), remote1)
		if err != nil {
			fmt.Println("Error?")
			return err
		}

	}
	return nil
}
