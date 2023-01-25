package progressbar

import (
	"fmt"
	"os"
	"time"
)

var minUpdateMs int64 = 1000
var spinner string = "/-|\\"

// ProgressBar keeps track of the progress of an operation and writes a progress meter to the terminal
type ProgressBar struct {
	maxPos     int64
	name       string
	lastMsgLen int
	closed     bool
	start      time.Time
	units      string
	curPos     int64
	spinPos    int
}

// NewProgressBar - create a new client instance
func NewProgressBar(length int64, name string, units string) *ProgressBar {
	return &ProgressBar{
		maxPos:     length,
		name:       name,
		closed:     false,
		lastMsgLen: -1,
		units:      units,
		curPos:     0,
		spinPos:    0,
	}
}

// Start - start the progress bar
func (pb *ProgressBar) Start() {
	pb.start = time.Now()
}

// Update - update the progress timer
func (pb *ProgressBar) Update(pos int64) {
	elapsed := time.Since(pb.start)

	if elapsed.Milliseconds() < minUpdateMs {
		// don't update too often
		return
	}

	termWidth := getTermCols()

	if pb.lastMsgLen > termWidth {
		pb.lastMsgLen = termWidth
	}

	pb.curPos = pos
	pb.spinPos++

	if pb.spinPos >= len(spinner) {
		pb.spinPos = 0
	}

	pct := float64(pb.curPos) / float64(pb.maxPos)

	remaining := pb.maxPos - pb.curPos
	rate := float64(pb.curPos) / float64(elapsed.Milliseconds()) // ticks per millisecond

	timeLeft := time.Duration(float64(remaining)/rate) * time.Millisecond // time left in milliseconds

	left := fmt.Sprintf("%s %s [", pb.name, string(spinner[pb.spinPos]))
	right := fmt.Sprintf("] %.1f%% %s ETA", pct*100, prettyDuration(timeLeft))

	barsize := termWidth - len(left) - len(right)

	bar := ""
	for i := 0; i < barsize; i++ {
		if i < int(float64(barsize)*pct) {
			bar += "="
		} else if i == int(float64(barsize)*pct) {
			bar += ">"
		} else {
			bar += "-"
		}
	}

	msg := left + bar + right

	pb.clearLine()
	os.Stderr.WriteString("\r" + msg)

	pb.lastMsgLen = len(msg)

}

// Done - complete the progress bar
func (pb *ProgressBar) Done() {
	if !pb.closed {
		pb.closed = true
		pb.clearLine()
		elapsed := time.Since(pb.start)

		os.Stderr.WriteString(fmt.Sprintf("\rElapsed time: %s\n", prettyDuration(elapsed)))
	}
}

// Error - complete the progress bar, but adds an error message
func (pb *ProgressBar) Error(msg string) {
	if !pb.closed {
		pb.closed = true
		pb.clearLine()

		os.Stderr.WriteString(msg)
	}
}

func (pb *ProgressBar) clearLine() {
	if pb.lastMsgLen > 0 {
		s := ""
		for i := 0; i < pb.lastMsgLen; i++ {
			s += " "
		}
		os.Stderr.WriteString("\r" + s)
	}
}

type tmpWinsize struct {
	rows uint16
	cols uint16
	xpix uint16
	ypix uint16
}

// func getTermCols() int {
// 	var ws tmpWinsize
// 	ret, _, _ := syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdout), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&ws)))
// 	if int(ret) == -1 {
// 		// if an error, set the width to a fixed 60 columns
// 		return 60
// 	}
// 	return int(ws.cols)
// }

func prettyDuration(d time.Duration) string {

	ret := ""
	sec := int(d.Seconds())

	if sec >= 3600 {
		h := sec / 3600
		sec = sec - (h * 3600)
		ret += fmt.Sprintf("%d:", h)
	}
	if sec >= 60 {
		m := sec / 60
		sec = sec - (m * 60)
		if m < 10 {
			ret += fmt.Sprintf("0%d:", m)
		} else {
			ret += fmt.Sprintf("%d:", m)
		}
	} else if ret != "" {
		ret += "00:"
	}

	if ret == "" {
		ret += fmt.Sprintf("%ds", sec)
	} else if sec < 10 {
		ret += fmt.Sprintf("0%d", sec)
	} else {
		ret += fmt.Sprintf("%d", sec)
	}

	return ret

}
