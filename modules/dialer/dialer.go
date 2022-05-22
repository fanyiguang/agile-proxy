package dialer

import (
	"net"
	"time"
)

type Dialer interface {
	Dial(network string, host, port string) (conn net.Conn, err error)
	DialTimeout(network string, host, port string, timeout time.Duration) (conn net.Conn, err error)
}
