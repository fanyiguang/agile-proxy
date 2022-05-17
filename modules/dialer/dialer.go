package dialer

import (
	"context"
	"net"
	"time"
)

type Dialer interface {
	Dial(network string, ip, port string) (conn net.Conn, err error)
	DialTimeout(network string, ip, port string, timeout time.Duration) (conn net.Conn, err error)
	DialContext(context context.Context, ip, port string, timeout time.Duration) (conn net.Conn, err error)
}
