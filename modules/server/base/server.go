package base

import (
	"agile-proxy/modules/ipc"
	"agile-proxy/modules/transport"
	"net"
)

type Server struct {
	ServerName  string
	ServerType  string
	Ip          string
	Port        string
	Username    string
	Password    string
	DoneCh      chan struct{}
	Listen      net.Listener
	Transmitter transport.Transport
	OutputMsgCh chan<- ipc.Msg
}

func (b *Server) Name() string {
	return b.ServerName
}

func (b *Server) Type() string {
	return b.ServerType
}
