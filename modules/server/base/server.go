package base

import (
	"agile-proxy/modules/plugin"
	"agile-proxy/modules/transport"
	"net"
)

type Server struct {
	plugin.Net
	plugin.Identity
	OutMsg      plugin.PipelineOutput
	DoneCh      chan struct{}
	Listen      net.Listener
	Transmitter transport.Transport
}
