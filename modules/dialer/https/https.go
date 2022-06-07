package https

import (
	"agile-proxy/modules/dialer/base"
	"agile-proxy/modules/plugin"
	"agile-proxy/pkg/https"
	"context"
	sysTls "crypto/tls"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"time"
)

type Https struct {
	base.Dialer
	plugin.Tls
	plugin.NetInfo
	httpsClient *https.Client
}

func (h *Https) Dial(network string, host, port string) (conn net.Conn, err error) {
	conn, err = h.BaseDial(network, h.Host, h.Port)
	if err != nil {
		return
	}

	config, err := h.CreateTlsConfig(h.Host)
	if err != nil {
		_ = conn.Close()
		return
	}

	tlsConn := sysTls.Client(conn, config)
	err = tlsConn.Handshake()
	if err != nil {
		_ = conn.Close()
		return
	}
	conn = tlsConn

	err = h.httpsClient.Handshake(conn, net.JoinHostPort(host, port))
	if err != nil {
		_ = conn.Close()
	}
	return
}

func (h *Https) DialTimeout(network string, host, port string, timeout time.Duration) (conn net.Conn, err error) {
	conn, err = h.BaseDialTimeout(network, h.Host, h.Port, timeout)
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

	err = h.httpsClient.Handshake(conn, net.JoinHostPort(host, port))
	if err != nil {
		_ = conn.Close()
	}
	return
}

func (h *Https) Close() (err error) {
	return
}

func New(jsonConfig json.RawMessage) (obj *Https, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "new")
		return
	}

	obj = &Https{
		Dialer: base.Dialer{
			IdentInfo: plugin.IdentInfo{
				ModuleName: config.Name,
				ModuleType: config.Type,
			},
			OutMsg: plugin.PipelineOutput{
				Ch: plugin.PipelineOutputCh,
			},
		},
		Tls: plugin.Tls{
			CrtPath: config.CrtPath,
			KeyPath: config.KeyPath,
		},
		NetInfo: plugin.NetInfo{
			Host:     config.Ip,
			Port:     config.Port,
			Username: config.Username,
			Password: config.Password,
		},
	}

	obj.httpsClient = https.New(config.Username, config.Password)

	return
}
