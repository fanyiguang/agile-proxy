package direct

import (
	"context"
	"net"
	"time"
)

type Direct struct {
}

func (d *Direct) Dial(network string, ip, port string) (conn net.Conn, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *Direct) DialTimeout(network string, ip, port string, timeout time.Duration) (conn net.Conn, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *Direct) DialContext(context context.Context, ip, port string, timeout time.Duration) (conn net.Conn, err error) {
	//TODO implement me
	panic("implement me")
}
