package socks5

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"nimble-proxy/modules/client"
	"nimble-proxy/modules/dialer"
	"nimble-proxy/pkg/socks5"
)

type Socks5 struct {
	client.BaseClient
	dialerName dialer.Dialer
	auth       int
}

func (s *Socks5) Dial(network string, host, port []byte) (conn net.Conn, err error) {
	s.dial(network)

	socks5Client := socks5.NewClient(conn, host, port, socks5.SetClientAuth(s.auth), socks5.SetClientUsername(s.Username), socks5.SetClientPassword(s.Password))
	err = socks5Client.HandShark()

	return
}

func (s *Socks5) Close() {
	//TODO implement me
	panic("implement me")
}

func (s *Socks5) dial(network string) (conn net.Conn, err error) {
	if s.dialerName != nil {
		conn, err = s.dialerName.Dial(network, s.Host, s.Port)
		if err != nil && s.Mode == 1 { // 连接器返回失败且自身为严格模式直接返回
			return
		}
	}

	conn, err = net.Dial(network, net.JoinHostPort(s.Host, s.Port))
	if err != nil {
		err = errors.Wrap(err, "socks5 Dial")
	}
	return
}

func New(strConfig string) (obj *Socks5, err error) {
	var config Config
	err = json.Unmarshal([]byte(strConfig), &config)
	if err != nil {
		err = errors.Wrap(err, "socks5 new")
		return
	}

	obj = &Socks5{
		BaseClient: client.BaseClient{
			Host:     config.Ip,
			Port:     config.Port,
			Username: config.Username,
			Password: config.Password,
		},
		auth: config.Auth,
	}

	return
}
