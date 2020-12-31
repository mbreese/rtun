package server

import (
	"strings"

	"github.com/mbreese/rtun/server/oshandler"
)

// NotifyCmd - echo back to the client
type NotifyCmd struct{}

func init() {
	addHandler("NOTIFY ", &NotifyCmd{})
}
func (cmd *NotifyCmd) destroy() {}

func (cmd *NotifyCmd) handle(ctx CmdContext) error {
	val := ctx.cmd[7:]
	val = strings.TrimSpace(val)
	val = strings.ReplaceAll(val, "$", "")
	val = strings.ReplaceAll(val, "\r", "")
	val = strings.ReplaceAll(val, "\n", "")
	val = strings.ReplaceAll(val, "\"", "")
	val = strings.ReplaceAll(val, "'", "")
	val = strings.ReplaceAll(val, ";", "")

	// fmt.Printf("Alert: %s\n", val)

	ar := strings.Split(val, " ")

	msg := strings.Join(ar[1:], " ")

	oshandler.Notify(msg, ar[0])
	_, err := ctx.write([]byte("OK\r\n"))
	return err
}
