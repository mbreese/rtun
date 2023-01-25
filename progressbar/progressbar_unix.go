// go:build linux freebsd darwin
// +build linux freebsd darwin !windows

package progressbar

import (
	"syscall"
	"unsafe"
)

func getTermCols() int {
	var ws tmpWinsize
	ret, _, _ := syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdout), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&ws)))
	if int(ret) == -1 {
		// if an error, set the width to a fixed 60 columns
		return 60
	}
	return int(ws.cols)
}
