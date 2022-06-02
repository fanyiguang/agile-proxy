package base

import (
	"agile-proxy/modules/base"
	"agile-proxy/modules/transport"
	"net"
)

type Server struct {
	base.NetInfo
	base.IdentInfo
	base.OutputMsg
	DoneCh      chan struct{}
	Listen      net.Listener
	Transmitter transport.Transport
}
