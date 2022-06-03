package https

import (
	commonBase "agile-proxy/modules/base"
	"agile-proxy/modules/client/base"
	"agile-proxy/pkg/https"
	"context"
	sysTls "crypto/tls"
	"net"
	"time"
)

type Https struct {
	base.Client
	commonBase.Tls
	httpsClient *https.Client
	tlsConfig   *sysTls.Config
}

func (h *Https) Dial(network string, host, port []byte) (conn net.Conn, err error) {
	conn, err = h.Dialer.Dial(network, h.Host, h.Port)
	if err != nil {
		return
	}

	config, err := h.CreateTlsConfig(h.Host)
	if err != nil {
		_ = conn.Close()
		return
	}

	conn, err = h.Handshake(context.Background(), conn, config)
	if err != nil {
		_ = conn.Close()
		return
	}

	err = h.httpsClient.Handshake(conn, net.JoinHostPort(string(host), h.GetStrPort(port)))
	if err != nil {
		_ = conn.Close()
	}
	return
}

func (h *Https) DialTimeout(network string, host, port []byte, timeout time.Duration) (conn net.Conn, err error) {
	conn, err = h.Dialer.DialTimeout(network, h.Host, h.Port, timeout)
	if err != nil {
		return
	}

	config, err := h.CreateTlsConfig(h.Host)
	if err != nil {
		_ = conn.Close()
		return
	}

	withTimeout, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()
	conn, err = h.Handshake(withTimeout, conn, config)
	if err != nil {
		_ = conn.Close()
		return
	}

	err = h.httpsClient.Handshake(conn, net.JoinHostPort(string(host), h.GetStrPort(port)))
	if err != nil {
		_ = conn.Close()
	}
	return
}

func (h *Https) Close() (err error) {
	return
}
