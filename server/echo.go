package server

// EchoCmd - echo back to the client
type EchoCmd struct{}

func init() {
	addHandler("ECHO ", &EchoCmd{})
}
func (cmd *EchoCmd) destroy() {}

func (cmd *EchoCmd) handle(ctx CmdContext) error {
	_, err := ctx.write([]byte("OK " + ctx.cmd[5:] + "\r\n"))
	return err
}
