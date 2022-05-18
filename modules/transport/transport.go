package transport

import "net"

var transports = make(map[string]Transport)

type Transport interface {
	Transport(conn net.Conn, ip, port string) (err error)
}

func GetTransport(name string) (t Transport) {
	return transports[name]
}

func GetAllTransports() map[string]Transport {
	return transports
}
