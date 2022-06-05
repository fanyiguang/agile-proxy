package base

import (
	"agile-proxy/modules/plugin"
	"agile-proxy/modules/transport"
	"net"
)

type Server struct {
	plugin.NetInfo
	plugin.IdentInfo
	OutMsg      plugin.PipelineOutput
	DoneCh      chan struct{}
	Listen      net.Listener
	Transmitter transport.Transport
}
