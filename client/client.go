package client

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"github.com/mbreese/rtun/progressbar"
)

var socketTimeout time.Duration = 120 * time.Second

// Client is the main connection to the server
type Client struct {
	socketFilename string
	conn           net.Conn
	reader         *bufio.Reader
	closed         bool
}

// // NewClient - create a new client instance
// func NewClient(fname string) *Client {
// 	return &Client{
// 		socketFilename: fname,
// 		conn:           nil,
// 		closed:         false,
// 	}
// }

// Connect - connect to the server
func Connect(fname string) *Client {

	// fmt.Printf("Connecting to: %s\n", fname)

	conn, err := net.Dial("unix", fname)
	if err != nil {
		panic(err)
	}

	return &Client{
		socketFilename: fname,
		conn:           conn,
		reader:         bufio.NewReader(conn),
		closed:         false,
	}
}

// Close client connection
func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) writeString(s string) (int, error) {
	return c.writeBytes([]byte(s))
}
func (c *Client) writeBytes(b []byte) (int, error) {
	// set a timeout of 120 seconds for read/write operations
	err := c.conn.SetDeadline(time.Now().Add(socketTimeout))
	if err != nil {
		panic(err)
	}
	return c.conn.Write(b)

}

func (c *Client) readLine() (string, error) {
	// set a timeout of 120 seconds for read/write operations
	err := c.conn.SetDeadline(time.Now().Add(socketTimeout))
	if err != nil {
		panic(err)
	}
	line, err := c.reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	line = strings.TrimRight(line, "\n")
	line = strings.TrimRight(line, "\r")

	return line, nil
}

// Ping server
func (c *Client) Ping() error {
	c.writeString("PING\r\n")

	line, err := c.readLine()

	if err != nil {
		return err
	}

	fmt.Println(line)
	return nil
}

// PID - get the PID of the server process
func (c *Client) PID() error {
	c.writeString("PID\r\n")

	line, err := c.readLine()

	if err != nil {
		return err
	}

	fmt.Println(line)
	return nil
}

// Shutdown server
func (c *Client) Shutdown() error {
	c.writeString("SHUTDOWN\r\n")

	line, err := c.readLine()

	if err != nil {
		return err
	}

	fmt.Println(line)
	return nil
}

// Echo server
func (c *Client) Echo(val string) error {
	c.writeString(fmt.Sprintf("ECHO %s\r\n", val))

	line, err := c.readLine()

	if err != nil {
		return err
	}

	fmt.Println(line)
	return nil
}

// Notify send a notification to the server
func (c *Client) Notify(val string) error {
	host, _ := os.Hostname()
	c.writeString(fmt.Sprintf("NOTIFY %s %s\r\n", host, val))

	line, err := c.readLine()

	if err != nil {
		return err
	}

	fmt.Println(line)
	return nil
}

// View a file on the server (sends the file, saving to a temp file, and then opens it)
func (c *Client) View(local string) error {
	return c.sendOrView(local, "", true)
}

// Send a file to the server, local and remote are full path names for files (not a directory)
func (c *Client) Send(local string, remote string) error {
	return c.sendOrView(local, remote, false)
}
func (c *Client) sendOrView(local string, remote string, isView bool) error {

	info, err := os.Stat(local)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("Local file cannot be a directory") // this should be taken care of in the cmd level
	}

	if remote == "" || remote == "." || remote == ".." || remote == "/" {
		remote = path.Base(local)
	}

	fmt.Printf("Uploading %s => %s\n", local, remote)
	if isView {
		c.writeString(fmt.Sprintf("VIEW %d %s\r\n", info.Size(), remote))
	} else {
		c.writeString(fmt.Sprintf("SEND %d %s\r\n", info.Size(), remote))
	}
	line, err := c.readLine()

	if err != nil {
		return err
	}

	if line != "READY" {
		return fmt.Errorf("Error writing file to server: %s", line)
	}

	// setup the progress bar

	pb := progressbar.NewProgressBar(info.Size(), local, "")

	f, err := os.Open(local)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()

	pb.Start()
	var total int = 0
	buf := make([]byte, 1024*1024)
	for {
		n, err := f.Read(buf)
		// fmt.Printf("Read %d bytes from file\n", n)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if n > 0 {
			total = total + n
			pb.Update(int64(total))
			h.Write(buf[:n])
			c.writeBytes(buf[:n])
			// fmt.Printf("Wrote %d bytes to server\n", n)
		}

		if total >= int(info.Size()) {
			break
		}
	}

	c.writeString(fmt.Sprintf("%x", h.Sum(nil)) + "\r\n")

	line, err = c.readLine()

	if err != nil {
		return err
	}
	if line[0:2] == "OK" {
		pb.Done()
	} else {
		pb.Error(line[3:])
	}
	return nil
}
