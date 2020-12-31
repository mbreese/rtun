package server

import (
	"container/list"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/mbreese/rtun/server/oshandler"
)

// SendCmd - save a file sent from the client
type SendCmd struct {
	tmpFiles *list.List
}

func init() {
	addHandler("SEND ", &SendCmd{tmpFiles: list.New()})
	addHandler("VIEW ", &SendCmd{tmpFiles: list.New()})
}

func (cmd *SendCmd) destroy() {
	for cmd.tmpFiles.Len() > 0 {
		head := cmd.tmpFiles.Front()
		fname := head.Value.(string)
		fmt.Printf("Removing temp/view file: %s", fname)
		_, err := os.Stat(fname)
		if err == nil {
			fmt.Println(", OK")
			os.Remove(fname)
		} else {
			fmt.Println(", missing")
		}
		cmd.tmpFiles.Remove(head)
	}
}

func (cmd *SendCmd) handle(ctx CmdContext) error {
	var sizeStr string
	var fname string

	isView := ctx.cmd[:5] == "VIEW "
	tmpFilename := ""

	for pos := 5; pos < len(ctx.cmd); pos++ {
		if ctx.cmd[pos] == ' ' {
			sizeStr = ctx.cmd[5:pos]
			fname = ctx.cmd[pos+1:]
			break
		}
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		ctx.write([]byte(fmt.Sprintf("ERROR 101 %s\r\n", err)))
		return err
	}

	var fout *os.File

	if isView {
		fout, err = ioutil.TempFile("", fmt.Sprintf("tmp_*_%s", path.Base(fname)))
		cmd.tmpFiles.PushBack(fout.Name())
		tmpFilename = fout.Name()
		fmt.Printf("Viewing file: %s (%d bytes) => %s\n", fname, size, tmpFilename)
	} else {
		fname = filepath.FromSlash(fname) // maybe help with Windows?

		for path.Clean(fname)[:3] == (".." + string(os.PathSeparator)) {
			fname = fname[3:]
		}

		fname = path.Join(ctx.server.saveDir, path.Clean(fname))

		finfo, err := os.Stat(fname)
		if err != nil && !os.IsNotExist(err) {
			// if it doesn't exist, that's okay
			ctx.write([]byte(fmt.Sprintf("ERROR  102 %s\r\n", err)))
			return err
		}

		if finfo != nil && finfo.IsDir() {
			// if it does exist, that's okay (we'll overwrite), unless it is a directory
			ctx.write([]byte(fmt.Sprintf("ERROR 103 Cannot write to file: %s\r\n", fname)))
			return errors.New("Target file is a directory")
		}

		_, err = os.Stat(path.Dir(fname))
		if err != nil {
			os.MkdirAll(path.Dir(fname), 0755)
		}

		fout, err = ioutil.TempFile(path.Dir(fname), "."+path.Base(fname))
		// f, err = os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			ctx.write([]byte(fmt.Sprintf("ERROR 104 %s\r\n", err)))
			return err
		}
		cmd.tmpFiles.PushBack(fout.Name())
	}

	fmt.Printf("Receiving file: %s (%d bytes)\n", fname, size)

	ctx.write([]byte("READY\r\n"))
	buf := make([]byte, 1024*1024)
	h := sha256.New()

	bytesToRead := len(buf)
	totalRead := 0

	for size-totalRead > 0 {
		if size-totalRead < bytesToRead {
			bytesToRead = size - totalRead
		}

		n, err := ctx.read(buf[:bytesToRead])
		// fmt.Printf("Read %d bytes from client (remaining: %d)\n", n, remaining)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		h.Write(buf[:n])
		fout.Write(buf[:n])

		totalRead += n
	}

	if err := fout.Close(); err != nil {
		ctx.write([]byte(fmt.Sprintf("ERROR 105 %s\r\n", err)))
		return err
	}

	if totalRead != size {
		ctx.write([]byte(fmt.Sprintf("ERROR 106 expected %d bytes, got %d bytes\r\n", size, totalRead)))
		return fmt.Errorf("Bad file size; Expected %d, Got %d", size, totalRead)
	}

	shasum, err := ctx.readLine()
	mysum := fmt.Sprintf("%x", h.Sum(nil))

	if shasum == mysum {
		if isView {
			oshandler.View(tmpFilename)
		} else {
			os.Rename(fout.Name(), fname)
		}
		fmt.Println("Success")

		ctx.write([]byte("OK\r\n"))
	} else {
		ctx.write([]byte(fmt.Sprintf("ERROR 107 expected %s, got %s\r\n", shasum, mysum)))
		return fmt.Errorf("Bad sha256; Expected %s, Got %s", shasum, mysum)
	}

	return nil
}
