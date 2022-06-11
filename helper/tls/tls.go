package tls

import (
	"agile-proxy/helper/log"
	"crypto/tls"
	"crypto/x509"
	"github.com/pkg/errors"
	"io/ioutil"
)

func CreateConfig(crtPath, keyPath, caPath string) (tlsConfig *tls.Config, err error) {
	if crtPath == "" || keyPath == "" { // 跳过认证
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		return
	}

	certificate, _err := tls.LoadX509KeyPair(crtPath, keyPath)
	if _err != nil {
		err = errors.Wrap(_err, "tls.LoadX509KeyPair")
		return
	}

	tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}
	if caPath != "" {
		pool, err := loadCa(caPath)
		if err != nil {
			log.WarnF("load ca failed: %v", err)
		} else {
			tlsConfig.RootCAs = pool
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}
	}
	return
}

func loadCa(caPath string) (cp *x509.CertPool, err error) {
	cp = x509.NewCertPool()
	data, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, err
	}
	if !cp.AppendCertsFromPEM(data) {
		return nil, errors.New("AppendCertsFromPEM failed")
	}
	return
}
