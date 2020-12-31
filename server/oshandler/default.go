// +build !darwin,!linux

package oshandler

import (
	"fmt"
	"runtime"
)

// Notify send a desktop notification
func Notify(msg string, title string) {
	fmt.Printf("No default NOTIFY handler for %s\n", runtime.GOOS)
}

// View a file
func View(fname string) {
	fmt.Printf("No default VIEW handler for %s\n", runtime.GOOS)
}
