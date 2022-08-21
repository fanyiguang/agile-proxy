package ssl

import (
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/client/base"
	"agile-proxy/pkg/socks5"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"time"
)

type Ssl struct {
	base.Client
	assembly.Tls
	socks5Client *socks5.Client
	authMode     int
}

func (s *Ssl) Dial(network string, host, port []byte) (conn net.Conn, err error) {
	conn, err = s.Client.Dial(network)
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

	err = s.socks5Client.HandShark(conn, host, port)
	if err != nil {
		_ = conn.Close()
	}
	return
}

func (s *Ssl) Run() (err error) {
	s.Client.Init()
	s.socks5Client = socks5.NewClient(socks5.SetClientAuth(s.authMode), socks5.SetClientUsername(s.Username), socks5.SetClientPassword(s.Password))
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
			Net:           assembly.CreateNet(config.Ip, config.Port, config.Username, config.Password),
			Identity:      assembly.CreateIdentity(config.Name, config.Type),
			Pipeline:      assembly.CreatePipeline(),
			PipelineInfos: config.PipelineInfos,
			Mode:          config.Mode,
			DialerName:    config.DialerName,
		},
		Tls:      assembly.CreateTls(config.CrtPath, config.KeyPath, config.CaPath, config.ServerName),
		authMode: config.AuthMode,
	}

	return
}
