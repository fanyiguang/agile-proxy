package socks5

import (
	commonBase "agile-proxy/modules/base"
	"agile-proxy/modules/client/base"
	"agile-proxy/modules/dialer"
	"agile-proxy/pkg/socks5"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"time"
)

type Socks5 struct {
	base.Client
	socks5Client *socks5.Client
	authMode     int
}

func (s *Socks5) Dial(network string, host, port []byte) (conn net.Conn, err error) {
	conn, err = s.Client.Dial(network)
	if err != nil {
		return
	}

	err = s.socks5Client.HandShark(conn, host, port)
	if err != nil {
		_ = conn.Close()
	}
	return
}

func (s *Socks5) DialTimeout(network string, host, port []byte, timeout time.Duration) (conn net.Conn, err error) {
	conn, err = s.Client.DialTimeout(network, timeout)
	if err != nil {
		return
	}

	err = s.socks5Client.HandShark(conn, host, port)
	if err != nil {
		_ = conn.Close()
	}
	return
}

func (s *Socks5) Close() (err error) {
	return
}

func New(strConfig json.RawMessage) (obj *Socks5, err error) {
	var config Config
	err = json.Unmarshal(strConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "new")
		return
	}

	obj = &Socks5{
		Client: base.Client{
			NetInfo: commonBase.NetInfo{
				Host:     config.Ip,
				Port:     config.Port,
				Username: config.Username,
				Password: config.Password,
			},
			IdentInfo: commonBase.IdentInfo{
				ModuleName: config.Name,
				ModuleType: config.Type,
			},
			OutputMsg: commonBase.OutputMsg{
				OutputMsgCh: commonBase.OutputCh,
			},
			Mode: config.Mode,
		},
		authMode: config.AuthMode,
	}

	if config.DialerName != "" {
		obj.Client.Dialer = dialer.GetDialer(config.DialerName)
	}
	obj.socks5Client = socks5.NewClient(socks5.SetClientAuth(obj.authMode), socks5.SetClientUsername(obj.Username), socks5.SetClientPassword(obj.Password))

	return
}
