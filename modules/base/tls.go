package base

import (
	"agile-proxy/helper/tls"
	"context"
	sysTls "crypto/tls"
	"net"
)

type Tls struct {
	TlsConfig *sysTls.Config
	CrtPath   string
	KeyPath   string
}

func (t *Tls) CreateTlsConfig(host string) (tlsConfig *sysTls.Config, err error) {
	if t.TlsConfig != nil {
		return t.TlsConfig, nil
	}

	tlsConfig, err = tls.CreateConfig(t.CrtPath, t.KeyPath)
	if err != nil {
		return
	}

	tlsConfig.ServerName = host
	t.TlsConfig = tlsConfig
	return
}

func (t *Tls) Handshake(ctx context.Context, rawConn net.Conn, config *sysTls.Config) (conn *sysTls.Conn, err error) {
	return tls.Handshake(ctx, rawConn, config)
}
