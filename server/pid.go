package server

import (
	"fmt"
	"os"
)

// PIDCmd - handle ping request
type PIDCmd struct{}

func init() {
	addHandler("PID", &PIDCmd{})
}

func (cmd *PIDCmd) destroy() {}

func (cmd *PIDCmd) handle(ctx CmdContext) error {
	_, err := ctx.write([]byte(fmt.Sprintf("OK PID %d, saveDir %s\r\n", os.Getpid(), ctx.server.saveDir)))
	return err
}
