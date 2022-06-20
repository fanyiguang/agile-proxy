package tls

import (
	"agile-proxy/helper/log"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
)

func CreateServerConfig(crtPath, keyPath, caPath string) (tlsConfig *tls.Config, err error) {
	var certificate tls.Certificate
	if crtPath == "" || keyPath == "" { // 使用默认证书
		certificate, err = tls.X509KeyPair(DefaultServerCrt(), DefaultServerKey())
	} else {
		certificate, err = tls.LoadX509KeyPair(crtPath, keyPath)
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("%v %v", crtPath, keyPath))
			return
		}
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

func CreateClientConfig(crtPath, keyPath, caPath, host string) (tlsConfig *tls.Config, err error) {
	if host == "" {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		return
	}

	var certificate tls.Certificate
	if crtPath == "" || keyPath == "" { // 使用默认证书
		certificate, err = tls.X509KeyPair(DefaultClientCrt(), DefaultClientKey())
	} else {
		certificate, err = tls.LoadX509KeyPair(crtPath, keyPath)
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("%v %v", crtPath, keyPath))
			return
		}
	}

	tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ServerName:   host,
	}

	if caPath != "" {
		pool, err := loadCa(caPath)
		if err != nil {
			log.WarnF("load ca failed: %v", err)
		} else {
			tlsConfig.RootCAs = pool
		}
	} else if crtPath == "" || keyPath == "" { // 走程序默认ca
		pool, err := loadDefaultCa()
		if err != nil {
			log.WarnF("load default ca failed: %v", err)
		} else {
			tlsConfig.RootCAs = pool
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

func loadDefaultCa() (cp *x509.CertPool, err error) {
	cp = x509.NewCertPool()
	if !cp.AppendCertsFromPEM(DefaultCaCrt()) {
		return nil, errors.New("AppendCertsFromPEM failed")
	}
	return
}
