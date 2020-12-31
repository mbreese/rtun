package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var socketTimeout time.Duration = 120 * time.Second

// Server is the main server that will echo input back to the client
type Server struct {
	socketFilename string
	saveDir        string
	listener       net.Listener
	closed         bool
	verbose        bool
}

// CmdContext - the request context
type CmdContext struct {
	server Server
	conn   net.Conn
	reader *bufio.Reader
	cmd    string
}

var serverHandlers = make(map[string]serverCmd)

func addHandler(k string, cmd serverCmd) {
	// fmt.Printf("Registering %s\n", k)
	serverHandlers[k] = cmd
}

// ServerCmd - interface to handle requests
type serverCmd interface {
	handle(cxt CmdContext) error
	destroy()
}

// NewServer - create a new echo server instance
func NewServer(socket string, saveDir string, verbose bool) *Server {
	return &Server{
		socketFilename: socket,
		saveDir:        saveDir,
		listener:       nil,
		closed:         false,
		verbose:        verbose,
	}
}

// Listen - start listening on the socket
func (s Server) Listen() {
	var wg sync.WaitGroup
	fmt.Printf("socket: %s, saveDir: %s, verbose: %t\n", s.socketFilename, s.saveDir, s.verbose)

	cTerm := make(chan os.Signal)
	signal.Notify(cTerm, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		for {
			sig := <-cTerm
			// fmt.Printf("Signal1: %s\n", sig)

			if sig == syscall.SIGHUP {
				// do nothing, this should in theory reload config, but not sure what to do with this program...
			} else {
				wg.Add(1) // the sig handler
				s.Close()
				wg.Done()
				// if it takes 10 seconds to close, force an exit
				time.Sleep(10 * time.Second)
				os.Exit(1)
			}
		}
	}()

	fmt.Printf("Listening on: %s\n", s.socketFilename)
	fmt.Printf("Saving files to: %s\n", s.saveDir)

	l, err := net.Listen("unix", s.socketFilename)
	if err != nil {
		panic(err)
	}
	if l == nil {
		panic("couldn't start listening: " + err.Error())
	}

	defer l.Close()
	s.listener = l

	os.Chmod(s.socketFilename, 0600)

	for !s.closed {
		conn, err := s.listener.Accept()

		if s.closed || conn == nil {
			break
		}

		if err != nil {
			panic(err)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err = conn.SetDeadline(time.Now().Add(socketTimeout))
			if err != nil {
				fmt.Printf("Error setting read/write timeouts.\n")
				return
			}
			buf := bufio.NewReader(conn)
			for {
				line, err := buf.ReadString('\n')
				if err != nil {
					if s.verbose {
						fmt.Println("Client disconnected.")
					}
					conn.Close()
					break
				}
				line = strings.TrimRight(line, "\n")
				line = strings.TrimRight(line, "\r")
				if s.verbose {
					fmt.Printf("< %s\n", line)
				}

				var ctx = CmdContext{server: s, conn: conn, reader: buf, cmd: line}

				var found = false
				for k, e := range serverHandlers {
					if strings.Index(line, k) == 0 {
						found = true
						err = e.handle(ctx)
					}
				}

				if !found {
					ctx.write([]byte(fmt.Sprintf("ERROR 100 Unknown command: %s", line)))
				}

				if err != nil {
					if s.verbose {
						fmt.Printf("Error processing request: %s.\n", err)
					}
					conn.Close()
					break
				}
			}
		}()
	}
	wg.Wait()

}

// Close - stop the server
func (s Server) Close() {
	if !s.closed {
		s.closed = true
		if s.listener != nil {
			// fmt.Printf("Shutting down listener\n")
			s.listener.Close()
			s.listener = nil
		}
		// fmt.Printf("Shutting down handlers\n")
		for _, e := range serverHandlers {
			// fmt.Printf(" => %v\n", e)
			e.destroy()
		}

		// this forces the socket to timeout
		net.Dial("unix", s.socketFilename)
	}
}

func (ctx CmdContext) readLine() (string, error) {
	err := ctx.conn.SetDeadline(time.Now().Add(socketTimeout))
	if err != nil {
		return "", err
	}
	line, err := ctx.reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	line = strings.TrimRight(line, "\n")
	line = strings.TrimRight(line, "\r")

	if ctx.server.verbose {
		fmt.Printf("< %s\n", line)
	}
	return line, nil
}

func (ctx CmdContext) write(b []byte) (int, error) {
	err := ctx.conn.SetDeadline(time.Now().Add(socketTimeout))
	if err != nil {
		return 0, err
	}
	if ctx.server.verbose {
		fmt.Printf("> %s", string(b))
	}
	return ctx.conn.Write(b)
}

func (ctx CmdContext) read(b []byte) (int, error) {
	err := ctx.conn.SetDeadline(time.Now().Add(socketTimeout))
	if err != nil {
		return 0, err
	}
	n, err := ctx.reader.Read(b)

	return n, err
}
