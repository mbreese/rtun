// go:build windows
package cmd

import "errors"

// Fork crete a new process
// see: https://github.com/immortal/immortal/blob/master/fork.go
func fork(stdout string) (int, error) {
	return -1, errors.New("-d doesn't work on Windows")
}
