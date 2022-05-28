package tls

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
)

func CreateConfig(crtPath, keyPath string) (tlsConfig *tls.Config, err error) {
	var bytes []byte
	pool := x509.NewCertPool()
	if crtPath == "" || keyPath == "" { // 忽略证书
		tlsConfig = &tls.Config{
			RootCAs: pool,
		}
		tlsConfig.InsecureSkipVerify = true
		return
	}

	bytes, err = ioutil.ReadFile(crtPath)
	if err != nil {
		err = errors.Wrap(err, "ReadFile")
		return
	}

	pool.AppendCertsFromPEM(bytes)
	tlsConfig = &tls.Config{
		RootCAs: pool,
	}
	var certificate tls.Certificate
	certificate, err = tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		err = errors.Wrap(err, "tls.LoadX509KeyPair")
		return
	}

	tlsConfig.Certificates = []tls.Certificate{certificate}
	return
}

func Handshake(ctx context.Context, rawConn net.Conn, config *tls.Config) (conn *tls.Conn, err error) {
	conn = tls.Client(rawConn, config)
	if err = conn.HandshakeContext(ctx); err != nil {
		rawConn.Close()
		return nil, errors.Wrap(err, "conn.HandshakeContext")
	}
	return conn, nil
}
