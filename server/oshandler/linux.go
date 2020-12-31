// +build linux

package oshandler

import (
	"fmt"
	"os/exec"
)

// Notify send a desktop notification
func Notify(msg string, title string) {
	cmd2 := exec.Command("notify-send", title, msg)
	err2 := cmd2.Run()
	if err2 != nil {
		fmt.Printf("Got err? %v (Notify on Linux requires \"notify-send\")\n", err2)
	}
}

// View a file
func View(fname string) {
	cmd2 := exec.Command("xdg-open", fname)
	err2 := cmd2.Run()
	if err2 != nil {
		fmt.Printf("Got err? %v (Notify on Linux requires \"xdg-open\")\n", err2)
	}
}
