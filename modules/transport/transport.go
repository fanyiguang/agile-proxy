package transport

import "net"

var transports = make(map[string]Transport)

type Transport interface {
	Transport(conn net.Conn, ip, port []byte) (err error)
}

type BaseTransport struct {
	Type       string
	Name       string
	ClientName string
	Mode       int
}

func GetTransport(name string) (t Transport) {
	return transports[name]
}

func GetAllTransports() map[string]Transport {
	return transports
}
