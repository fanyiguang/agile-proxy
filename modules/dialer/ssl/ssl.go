package ssl

import (
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/dialer/base"
	"agile-proxy/proxy/socks5"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"time"
)

type ssl struct {
	base.Dialer
	assembly.Tls
	assembly.Net
	socks5Client *socks5.Client
	authMode     int
}

func (s *ssl) Dial(network string, host, port string) (conn net.Conn, err error) {
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

func (s *ssl) DialTimeout(network string, host, port string, timeout time.Duration) (conn net.Conn, err error) {
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

func (s *ssl) Run() (err error) {
	err = s.init()
	return
}

func (s *ssl) Close() (err error) {
	return
}

func (s *ssl) init() (err error) {
	s.socks5Client = socks5.NewClient(socks5.SetClientAuth(s.authMode), socks5.SetClientUsername(s.Username), socks5.SetClientPassword(s.Password))
	return
}

func New(jsonConfig json.RawMessage) (obj *ssl, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "new")
		return
	}

	obj = &ssl{
		Dialer: base.Dialer{
			Net:           assembly.CreateNet(config.Ip, config.Port, config.Username, config.Password),
			Identity:      assembly.CreateIdentity(config.Name, config.Type),
			Pipeline:      assembly.CreatePipeline(),
			PipelineInfos: config.PipelineInfos,
			IFace:         config.Interface,
		},
		Tls:      assembly.CreateTls(config.CrtPath, config.KeyPath, config.CaPath, config.ServerName),
		authMode: config.AuthMode,
	}

	return
}
