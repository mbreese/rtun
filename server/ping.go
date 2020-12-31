package server

// PingCmd - handle ping request
type PingCmd struct{}

func init() {
	addHandler("PING", &PingCmd{})
}

func (cmd *PingCmd) destroy() {}

func (cmd *PingCmd) handle(ctx CmdContext) error {
	_, err := ctx.write([]byte("OK PONG\r\n"))
	return err
}
