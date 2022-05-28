package ssl

import (
	"context"
	sysTls "crypto/tls"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"nimble-proxy/helper/tls"
	"nimble-proxy/modules/client/base"
	"nimble-proxy/modules/dialer"
	"nimble-proxy/pkg/socks5"
	"time"
)

type Ssl struct {
	base.Client
	tlsConfig    *sysTls.Config
	socks5Client *socks5.Client
	crtPath      string
	keyPath      string
	authMode     int
}

func (s *Ssl) Dial(network string, host, port []byte) (conn net.Conn, err error) {
	conn, err = s.Client.Dial(network)
	if err != nil {
		return
	}

	config, err := s.createTlsConfig()
	if err != nil {
		return
	}

	conn, err = tls.Handshake(context.Background(), conn, config)
	if err != nil {
		return
	}

	err = s.socks5Client.HandShark(conn, host, port)
	return
}

func (s *Ssl) DialTimeout(network string, host, port []byte, timeout time.Duration) (conn net.Conn, err error) {
	conn, err = s.Client.DialTimeout(network, timeout)
	if err != nil {
		return
	}

	config, err := s.createTlsConfig()
	if err != nil {
		return
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()
	conn, err = tls.Handshake(ctx, conn, config)
	if err != nil {
		return
	}

	err = s.socks5Client.HandShark(conn, host, port)
	return
}

func (s *Ssl) Close() (err error) {
	//TODO 一些资源的释放减轻GC工作量
	return
}

func (s *Ssl) createTlsConfig() (tlsConfig *sysTls.Config, err error) {
	if s.tlsConfig != nil {
		return s.tlsConfig, nil
	}

	tlsConfig, err = tls.CreateConfig(s.crtPath, s.keyPath)
	if err != nil {
		return
	}

	tlsConfig.ServerName = s.Host
	s.tlsConfig = tlsConfig
	return
}

func New(strConfig json.RawMessage) (obj *Ssl, err error) {
	var config Config
	err = json.Unmarshal(strConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "socks5 new")
		return
	}

	obj = &Ssl{
		Client: base.Client{
			Host:       config.Ip,
			Port:       config.Port,
			Username:   config.Username,
			Password:   config.Password,
			ClientName: config.Name,
			ClientType: config.Type,
			Mode:       config.Mode,
		},
		crtPath:  config.CrtPath,
		keyPath:  config.KeyPath,
		authMode: config.AuthMode,
	}

	if config.DialerName != "" {
		obj.Client.Dialer = dialer.GetDialer(config.DialerName)
	}
	obj.socks5Client = socks5.NewClient(socks5.SetClientAuth(obj.authMode), socks5.SetClientUsername(obj.Username), socks5.SetClientPassword(obj.Password))

	return
}
