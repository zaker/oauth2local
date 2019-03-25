package ipc

import (
	"net"

	winio "github.com/Microsoft/go-winio"
)

type IPCServer struct {
	listener net.Listener
}

const pipeName = `\\.\pipe\oauth2local`

func Listen() (*IPCServer, error) {
	l, err := winio.ListenPipe(pipeName, nil)
	return &IPCServer{listener: l}, err
}
