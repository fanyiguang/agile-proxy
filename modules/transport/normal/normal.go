package normal

import (
	"net"
	"nimble-proxy/modules/transport"
)

type Normal struct {
	transport.BaseTransport
}

func (n Normal) Transport(conn net.Conn, ip, port []byte) (err error) {
	return
}
