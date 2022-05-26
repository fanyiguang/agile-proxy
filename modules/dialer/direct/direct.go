package direct

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"nimble-proxy/modules/dialer/base"
	"time"
)

type Direct struct {
	base.Dialer
}

func (d *Direct) Dial(network string, host, port string) (conn net.Conn, err error) {
	if d.IFace != "" {
		rAddr, _err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, port))
		if _err != nil {
			err = errors.Wrap(_err, "net.ResolveTCPAddr-1")
			return
		}
		lAddr, _err := net.ResolveTCPAddr("tcp", net.JoinHostPort(d.IFace, "0"))
		if _err != nil {
			err = errors.Wrap(_err, "net.ResolveTCPAddr-2")
			return
		}

		conn, err = net.DialTCP(network, lAddr, rAddr)
		if err != nil {
			err = errors.Wrap(err, "net.DialTCP")
		}
		return
	}

	return net.Dial(network, net.JoinHostPort(host, port))
}

func (d *Direct) DialTimeout(network string, ip, port string, timeout time.Duration) (conn net.Conn, err error) {
	//TODO implement me
	panic("implement me")
}

func New(jsonConfig json.RawMessage) (obj *Direct, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "direct new")
		return
	}

	obj = &Direct{
		Dialer: base.Dialer{
			DialerName: config.Name,
			DialerType: config.Type,
			IFace:      config.Interface,
		},
	}

	return
}
