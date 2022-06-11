package https

import (
	"agile-proxy/helper/common"
	"agile-proxy/modules/client/base"
	"agile-proxy/modules/dialer"
	"agile-proxy/modules/plugin"
	"agile-proxy/pkg/https"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"time"
)

type Https struct {
	base.Client
	plugin.Tls
	httpsClient *https.Client
}

func (h *Https) Dial(network string, host, port []byte) (conn net.Conn, err error) {
	conn, err = h.Client.Dial(network)
	if err != nil {
		return
	}

	config, err := h.CreateTlsConfig()
	if err != nil {
		_ = conn.Close()
		return
	}

	conn, err = h.Handshake(context.Background(), conn, config)
	if err != nil {
		_ = conn.Close()
		return
	}

	err = h.httpsClient.Handshake(conn, net.JoinHostPort(common.BytesToStr(host), h.GetStrPort(port)))
	if err != nil {
		_ = conn.Close()
	}
	return
}

func (h *Https) DialTimeout(network string, host, port []byte, timeout time.Duration) (conn net.Conn, err error) {
	conn, err = h.Client.DialTimeout(network, timeout)
	if err != nil {
		return
	}

	config, err := h.CreateTlsConfig()
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

	err = h.httpsClient.Handshake(conn, net.JoinHostPort(common.BytesToStr(host), h.GetStrPort(port)))
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
		Client: base.Client{
			Net: plugin.Net{
				Host:     config.Ip,
				Port:     config.Port,
				Username: config.Username,
				Password: config.Password,
			},
			Identity: plugin.Identity{
				ModuleName: config.Name,
				ModuleType: config.Type,
			},
			OutMsg: plugin.PipelineOutput{
				Ch: plugin.PipelineOutputCh,
			},
			Mode: config.Mode,
		},
		Tls: plugin.Tls{
			CrtPath: config.CrtPath,
			KeyPath: config.KeyPath,
		},
	}

	if config.DialerName != "" {
		obj.Client.Dialer = dialer.GetDialer(config.DialerName)
	}
	obj.httpsClient = https.New(config.Username, config.Password)

	return
}
