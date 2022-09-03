package ssl

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"time"
)

type DialSsl struct {
}

func (s *DialSsl) Dial(network, addr string) (conn net.Conn, err error) {
	dialer := &net.Dialer{}
	conn, err = tls.DialWithDialer(dialer, network, addr, getTlsConfig())
	return
}

func (s *DialSsl) DialContext(ctx context.Context, network, address string) (conn net.Conn, err error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(30 * time.Second)
	}

	dialer := &net.Dialer{
		Deadline: deadline,
	}
	conn, err = tls.DialWithDialer(dialer, network, address, getTlsConfig())
	return
}

/*获取tls配置*/
func getTlsConfig() (conf *tls.Config) {
	pool := x509.NewCertPool()
	conf = &tls.Config{
		RootCAs: pool,
	}
	/*证书不存在则不使用证书*/
	conf.InsecureSkipVerify = true
	return
}
