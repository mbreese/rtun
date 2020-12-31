package server

import (
	"fmt"
)

// ShutdownCmd - handle quit
type ShutdownCmd struct{}

func init() {
	addHandler("SHUTDOWN", &ShutdownCmd{})
}

func (cmd *ShutdownCmd) destroy() {}

func (cmd *ShutdownCmd) handle(ctx CmdContext) error {
	fmt.Println("Got shutdown command")
	ctx.write([]byte("OK Shutting down server...\r\n"))
	ctx.conn.Close()

	// time.Sleep(2 * time.Second)
	ctx.server.Close()
	return nil
}
