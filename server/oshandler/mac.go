// +build darwin

package oshandler

import (
	"fmt"
	"os/exec"
)

// Notify send a desktop notification
func Notify(msg string, title string) {
	cmd2 := exec.Command("osascript", "-e", fmt.Sprintf("display notification \"%s\" with title \"%s\"", msg, title))
	err2 := cmd2.Run()
	if err2 != nil {
		fmt.Printf("Got err? %v\n", err2)
	}
}

// View a file
func View(fname string) {
	cmd2 := exec.Command("open", fname)
	err2 := cmd2.Run()
	if err2 != nil {
		fmt.Printf("Got err? %v\n", err2)
	}
}
