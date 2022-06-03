package http

import (
	"agile-proxy/modules/client/base"
	"agile-proxy/pkg/https"
	sysTls "crypto/tls"
	"net"
	"time"
)

type Http struct {
	base.Client
	httpsClient *https.Client
	tlsConfig   *sysTls.Config
}

func (h *Http) Dial(network string, host, port []byte) (conn net.Conn, err error) {
	conn, err = h.Dialer.Dial(network, h.Host, h.Port)
	if err != nil {
		return
	}

	err = h.httpsClient.Handshake(conn, net.JoinHostPort(string(host), h.GetStrPort(port)))
	if err != nil {
		_ = conn.Close()
	}
	return
}

func (h *Http) DialTimeout(network string, host, port []byte, timeout time.Duration) (conn net.Conn, err error) {
	conn, err = h.Dialer.DialTimeout(network, h.Host, h.Port, timeout)
	if err != nil {
		return
	}

	err = h.httpsClient.Handshake(conn, net.JoinHostPort(string(host), h.GetStrPort(port)))
	if err != nil {
		_ = conn.Close()
	}
	return
}

func (h *Http) Close() (err error) {
	return
}
