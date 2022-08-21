package assembly

import (
	"agile-proxy/helper/tls"
	"context"
	sysTls "crypto/tls"
	"github.com/pkg/errors"
	"net"
)

type Tls struct {
	TlsConfig  *sysTls.Config
	CrtPath    string
	KeyPath    string
	CaPath     string
	ServerName string
}

func (t *Tls) CreateServerTlsConfig() (tlsConfig *sysTls.Config, err error) {
	if t.TlsConfig != nil {
		return t.TlsConfig, nil
	}

	tlsConfig, err = tls.CreateServerConfig(t.CrtPath, t.KeyPath, t.CaPath)
	if err != nil {
		return
	}

	t.TlsConfig = tlsConfig
	return
}

func (t *Tls) CreateClientTlsConfig() (tlsConfig *sysTls.Config, err error) {
	if t.TlsConfig != nil {
		return t.TlsConfig, nil
	}

	tlsConfig, err = tls.CreateClientConfig(t.CrtPath, t.KeyPath, t.CaPath, t.ServerName)
	if err != nil {
		return
	}

	t.TlsConfig = tlsConfig
	return
}

func (t *Tls) Handshake(ctx context.Context, rawConn net.Conn, config *sysTls.Config) (conn *sysTls.Conn, err error) {
	conn = sysTls.Client(rawConn, config)
	err = conn.HandshakeContext(ctx)
	if err != nil {
		_ = conn.Close()
		err = errors.Wrap(err, "Handshake.HandshakeContext")
		return
	}
	return
}

func CreateTls(crtPath, keyPath, caPath, serverName string) Tls {
	return Tls{
		CrtPath:    crtPath,
		KeyPath:    keyPath,
		CaPath:     caPath,
		ServerName: serverName,
	}
}
