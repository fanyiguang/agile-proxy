package ssh

import (
	sysSsh "golang.org/x/crypto/ssh"
	"net"
	"time"
)

type Conn struct {
	sysSsh.Channel
	localAddr  net.Addr
	remoteAddr net.Addr
}

func (c *Conn) LocalAddr() net.Addr {
	return c.localAddr
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *Conn) SetDeadline(t time.Time) error {
	return nil
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return nil
}
