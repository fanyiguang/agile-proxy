package socks5

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"nimble-proxy/helper/log"
	"nimble-proxy/modules/client/base"
	"nimble-proxy/modules/dialer"
	"nimble-proxy/pkg/socks5"
)

type Socks5 struct {
	base.Client
	dialer   dialer.Dialer
	authMode int
}

func (s *Socks5) Dial(network string, host, port []byte) (conn net.Conn, err error) {
	conn, err = s.dial(network)
	if err != nil {
		return
	}

	socks5Client := socks5.NewClient(conn, host, port, socks5.SetClientAuth(s.authMode), socks5.SetClientUsername(s.Username), socks5.SetClientPassword(s.Password))
	err = socks5Client.HandShark()
	return
}

func (s *Socks5) Close() {
	//TODO implement me
	panic("implement me")
}

func (s *Socks5) dial(network string) (conn net.Conn, err error) {
	if s.dialer != nil {
		conn, err = s.dialer.Dial(network, s.Host, s.Port)
		if err == nil || s.Mode == 1 { // mode=1 严格模式
			return
		}

		if err != nil {
			log.WarnF("s.dialer.Dial failed: %v", err)
		}
	}

	conn, err = net.Dial(network, net.JoinHostPort(s.Host, s.Port))
	if err != nil {
		err = errors.Wrap(err, "socks5 Dial")
	}
	return
}

func New(strConfig json.RawMessage) (obj *Socks5, err error) {
	var config Config
	err = json.Unmarshal(strConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "socks5 new")
		return
	}

	obj = &Socks5{
		Client: base.Client{
			Host:       config.Ip,
			Port:       config.Port,
			Username:   config.Username,
			Password:   config.Password,
			ClientName: config.Name,
			ClientType: config.Type,
			Mode:       config.Mode,
		},
		authMode: config.AuthMode,
	}

	if config.DialerName != "" {
		obj.dialer = dialer.GetDialer(config.DialerName)
	}

	return
}
