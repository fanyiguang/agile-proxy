package socks5

import (
	"agile-proxy/helper/common"
	"agile-proxy/modules/dialer/base"
	"agile-proxy/modules/plugin"
	"agile-proxy/pkg/socks5"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"time"
)

type Socks5 struct {
	base.Dialer
	plugin.Net
	socks5Client *socks5.Client
	authMode     int
}

func (s *Socks5) Dial(network string, host, port string) (conn net.Conn, err error) {
	conn, err = s.BaseDial(network, s.Host, s.Port)
	if err != nil {
		return
	}

	err = s.socks5Client.HandShark(conn, common.StrToBytes(host), common.StrToBytes(port))
	if err != nil {
		_ = conn.Close()
	}
	return
}

func (s *Socks5) DialTimeout(network string, host, port string, timeout time.Duration) (conn net.Conn, err error) {
	conn, err = s.BaseDialTimeout(network, s.Host, s.Port, timeout)
	if err != nil {
		return
	}

	err = s.socks5Client.HandShark(conn, common.StrToBytes(host), common.StrToBytes(port))
	if err != nil {
		_ = conn.Close()
	}
	return
}

func (s *Socks5) Close() (err error) {
	return
}

func New(jsonConfig json.RawMessage) (obj *Socks5, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "new")
		return
	}

	obj = &Socks5{
		Dialer: base.Dialer{
			Identity: plugin.Identity{
				ModuleName: config.Name,
				ModuleType: config.Type,
			},
			OutMsg: plugin.PipelineOutput{
				Ch: plugin.PipelineOutputCh,
			},
		},
		Net: plugin.Net{
			Host:     config.Ip,
			Port:     config.Port,
			Username: config.Username,
			Password: config.Password,
		},
		authMode: config.AuthMode,
	}

	obj.socks5Client = socks5.NewClient(socks5.SetClientAuth(obj.authMode), socks5.SetClientUsername(obj.Username), socks5.SetClientPassword(obj.Password))

	return
}
