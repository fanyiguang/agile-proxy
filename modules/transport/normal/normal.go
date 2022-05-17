package normal

import (
	"net"
)

type Normal struct {
}

func (n Normal) Transport(conn net.Conn, ip, port string) (err error) {
	//TODO implement me
	panic("implement me")
}
