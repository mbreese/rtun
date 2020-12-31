package main

import (
	"github.com/mbreese/rtun/cmd"
)

func main() {
	// if os.Args[1] == "server" {
	// 	args := make([]string, len(os.Args)-1)
	// 	forkme := false
	// 	for i, arg := range os.Args[1:] {
	// 		if arg == "-d" || arg == "--daemon" {
	// 			forkme = true
	// 		} else if forkme {
	// 			args[i-1] = arg
	// 		} else {
	// 			args[i] = arg
	// 		}
	// 	}
	// 	if forkme {
	// 		pid, _ := Fork(args, "test.out")
	// 		fmt.Printf("%d\n", pid)
	// 		os.Exit(0)
	// 	}

	// }

	cmd.Execute()
}
