package direct

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"nimble-proxy/modules/dialer"
	"time"
)

type Direct struct {
	dialer.BaseDialer
}

func (d *Direct) Dial(network string, host, port string) (conn net.Conn, err error) {
	if d.IFace != "" {
		rAddr, _err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, port))
		if _err != nil {
			err = errors.Wrap(_err, "net.ResolveTCPAddr")
			return
		}
		lAddr, _err := net.ResolveTCPAddr("tcp", net.JoinHostPort(d.IFace, ":0"))
		if _err != nil {
			err = errors.Wrap(_err, "net.ResolveTCPAddr")
			return
		}

		conn, err = net.DialTCP(network, lAddr, rAddr)
		return
	}

	return net.Dial(network, net.JoinHostPort(host, port))
}

func (d *Direct) DialTimeout(network string, ip, port string, timeout time.Duration) (conn net.Conn, err error) {
	//TODO implement me
	panic("implement me")
}

func New(jsonConfig string) (obj *Direct, err error) {
	var config Config
	err = json.Unmarshal([]byte(jsonConfig), &config)
	if err != nil {
		err = errors.Wrap(err, "direct new")
		return
	}

	obj = &Direct{
		dialer.BaseDialer{
			Name:  config.Name,
			Type:  config.Type,
			IFace: config.Interface,
		},
	}

	return
}
