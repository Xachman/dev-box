package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type ConsoleController struct {
}

func hijack(addr, method, path string, setRawTerminal bool, in io.ReadCloser, stdout, stderr io.Writer, started chan io.Closer, data interface{}) error {

	params := bytes.NewBufferString("{\"Detach\": false, \"Tty\": true}")
	req, err := http.NewRequest(method, path, params)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Docker-Client")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "tcp")
	req.Host = addr

	dial, err := net.Dial("tcp", addr)
	// When we set up a TCP connection for hijack, there could be long periods
	// of inactivity (a long running command with no output) that in certain
	// network setups may cause ECONNTIMEOUT, leaving the client in an unknown
	// state. Setting TCP KeepAlive on the socket connection will prohibit
	// ECONNTIMEOUT unless the socket connection truly is broken
	if tcpConn, ok := dial.(*net.TCPConn); ok {
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)
	}
	if err != nil {
		return err
	}
	clientconn := httputil.NewClientConn(dial, nil)
	defer clientconn.Close()

	// Server hijacks the connection, error 'connection closed' expected
	clientconn.Do(req)

	rwc, br := clientconn.Hijack()
	defer rwc.Close()

	if started != nil {
		started <- rwc
	}

	var receiveStdout chan error

	if stdout != nil || stderr != nil {
		go func() (err error) {
			if setRawTerminal && stdout != nil {
				_, err = io.Copy(stdout, br)
			}
			return err
		}()
	}

	go func() error {
		if in != nil {
			io.Copy(rwc, in)
		}

		if conn, ok := rwc.(interface {
			CloseWrite() error
		}); ok {
			if err := conn.CloseWrite(); err != nil {
			}
		}
		return nil
	}()

	if stdout != nil || stderr != nil {
		if err := <-receiveStdout; err != nil {
			return err
		}
	}
	spew.Dump(br)
	go func() {
		for {
			fmt.Println(br)
			spew.Dump(br)
		}
	}()

	return nil
}
