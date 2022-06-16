package ssl

import (
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/modules/dialer/base"
	"agile-proxy/modules/plugin"
	"agile-proxy/pkg/socks5"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"time"
)

type Ssl struct {
	base.Dialer
	plugin.Tls
	plugin.Net
	socks5Client *socks5.Client
	authMode     int
}

func (s *Ssl) Dial(network string, host, port string) (conn net.Conn, err error) {
	conn, err = s.BaseDial(network, s.Host, s.Port)
	if err != nil {
		return
	}

	config, err := s.CreateClientTlsConfig()
	if err != nil {
		_ = conn.Close()
		return
	}

	conn, err = s.Handshake(context.Background(), conn, config)
	if err != nil {
		_ = conn.Close()
		return
	}

	err = s.socks5Client.HandShark(conn, common.StrToBytes(host), common.StrToBytes(port))
	if err != nil {
		_ = conn.Close()
	}
	log.DebugF("ssl dialer link status: %v %v", err, net.JoinHostPort(host, port))
	return
}

func (s *Ssl) DialTimeout(network string, host, port string, timeout time.Duration) (conn net.Conn, err error) {
	conn, err = s.BaseDialTimeout(network, s.Host, s.Port, timeout)
	if err != nil {
		return
	}

	config, err := s.CreateClientTlsConfig()
	if err != nil {
		_ = conn.Close()
		return
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	conn, err = s.Handshake(ctx, conn, config)
	if err != nil {
		_ = conn.Close()
		return
	}

	err = s.socks5Client.HandShark(conn, common.StrToBytes(host), common.StrToBytes(port))
	if err != nil {
		_ = conn.Close()
	}
	return
}

func (s *Ssl) Close() (err error) {
	return
}

func New(jsonConfig json.RawMessage) (obj *Ssl, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "new")
		return
	}

	obj = &Ssl{
		Dialer: base.Dialer{
			Identity: plugin.Identity{
				ModuleName: config.Name,
				ModuleType: config.Type,
			},
			OutMsg: plugin.PipelineOutput{
				Ch: plugin.PipelineOutputCh,
			},
		},
		Tls: plugin.Tls{
			CrtPath:    config.CrtPath,
			KeyPath:    config.KeyPath,
			CaPath:     config.CaPath,
			ServerName: config.ServerName,
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
