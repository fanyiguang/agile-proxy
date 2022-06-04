package base

import (
	"agile-proxy/modules/plugin"
	"github.com/pkg/errors"
	"net"
	"time"
)

type Dialer struct {
	plugin.IdentInfo
	plugin.OutputMsg
	IFace string
}

func (d *Dialer) BaseDial(network string, host, port string) (conn net.Conn, err error) {
	if d.IFace != "" {
		return d.DialByIFace(network, host, port)
	}

	return net.Dial(network, net.JoinHostPort(host, port))
}

func (d *Dialer) BaseDialTimeout(network string, host, port string, timeout time.Duration) (conn net.Conn, err error) {
	if d.IFace != "" {
		return d.DialByIFace(network, host, port)
	}

	return net.DialTimeout(network, net.JoinHostPort(host, port), timeout)
}

func (d *Dialer) DialByIFace(network, host, port string) (conn net.Conn, err error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
		rAddr, _err := net.ResolveTCPAddr(network, net.JoinHostPort(host, port))
		if _err != nil {
			err = errors.Wrap(_err, "net.ResolveTCPAddr-1")
			return
		}
		lAddr, _err := net.ResolveTCPAddr(network, net.JoinHostPort(d.IFace, "0"))
		if _err != nil {
			err = errors.Wrap(_err, "net.ResolveTCPAddr-2")
			return
		}
		conn, err = net.DialTCP(network, lAddr, rAddr)
	case "udp", "udp4", "udp6":
		rAddr, _err := net.ResolveUDPAddr(network, net.JoinHostPort(host, port))
		if _err != nil {
			err = errors.Wrap(_err, "net.ResolveTCPAddr-1")
			return
		}
		lAddr, _err := net.ResolveUDPAddr(network, net.JoinHostPort(d.IFace, "0"))
		if _err != nil {
			err = errors.Wrap(_err, "net.ResolveTCPAddr-2")
			return
		}
		conn, err = net.DialUDP(network, lAddr, rAddr)
	}
	if err != nil {
		err = errors.Wrap(err, "net.DialTCP")
	}
	return
}
