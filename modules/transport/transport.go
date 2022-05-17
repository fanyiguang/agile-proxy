package transport

import "net"

type Transport interface {
	Transport(conn net.Conn, ip, port string) (err error)
}
