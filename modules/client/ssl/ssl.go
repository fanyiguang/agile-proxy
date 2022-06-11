package ssl

import (
	"agile-proxy/modules/client/base"
	"agile-proxy/modules/dialer"
	"agile-proxy/modules/plugin"
	"agile-proxy/pkg/socks5"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"time"
)

type Ssl struct {
	base.Client
	plugin.Tls
	socks5Client *socks5.Client
	authMode     int
}

func (s *Ssl) Dial(network string, host, port []byte) (conn net.Conn, err error) {
	conn, err = s.Client.Dial(network)
	if err != nil {
		return
	}

	config, err := s.CreateTlsConfig()
	if err != nil {
		_ = conn.Close()
		return
	}

	conn, err = s.Handshake(context.Background(), conn, config)
	if err != nil {
		_ = conn.Close()
		return
	}

	err = s.socks5Client.HandShark(conn, host, port)
	if err != nil {
		_ = conn.Close()
	}
	return
}

func (s *Ssl) DialTimeout(network string, host, port []byte, timeout time.Duration) (conn net.Conn, err error) {
	conn, err = s.Client.DialTimeout(network, timeout)
	if err != nil {
		return
	}

	config, err := s.CreateTlsConfig()
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

	err = s.socks5Client.HandShark(conn, host, port)
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
		Client: base.Client{
			Net: plugin.Net{
				Host:     config.Ip,
				Port:     config.Port,
				Username: config.Username,
				Password: config.Password,
			},
			Identity: plugin.Identity{
				ModuleName: config.Name,
				ModuleType: config.Type,
			},
			OutMsg: plugin.PipelineOutput{
				Ch: plugin.PipelineOutputCh,
			},
			Mode: config.Mode,
		},
		Tls: plugin.Tls{
			CrtPath: config.CrtPath,
			KeyPath: config.KeyPath,
		},
		authMode: config.AuthMode,
	}

	if config.DialerName != "" {
		obj.Client.Dialer = dialer.GetDialer(config.DialerName)
	}
	obj.socks5Client = socks5.NewClient(socks5.SetClientAuth(obj.authMode), socks5.SetClientUsername(obj.Username), socks5.SetClientPassword(obj.Password))

	return
}
