package base

import (
	"net"
	"nimble-proxy/modules/ipc"
	"nimble-proxy/modules/transport"
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