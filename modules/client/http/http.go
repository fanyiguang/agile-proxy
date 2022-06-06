package http

import (
	"agile-proxy/modules/client/base"
	"agile-proxy/modules/dialer"
	"agile-proxy/modules/plugin"
	"agile-proxy/pkg/https"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"time"
)

type Http struct {
	base.Client
	httpsClient *https.Client
}

func (h *Http) Dial(network string, host, port []byte) (conn net.Conn, err error) {
	conn, err = h.Client.Dial(network)
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
	conn, err = h.Client.DialTimeout(network, timeout)
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

func New(jsonConfig json.RawMessage) (obj *Http, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "new")
		return
	}

	obj = &Http{
		Client: base.Client{
			NetInfo: plugin.NetInfo{
				Host:     config.Ip,
				Port:     config.Port,
				Username: config.Username,
				Password: config.Password,
			},
			IdentInfo: plugin.IdentInfo{
				ModuleName: config.Name,
				ModuleType: config.Type,
			},
			OutMsg: plugin.PipelineOutput{
				Ch: plugin.PipelineOutputCh,
			},
			Mode: config.Mode,
		},
	}

	if config.DialerName != "" {
		obj.Client.Dialer = dialer.GetDialer(config.DialerName)
	}
	obj.httpsClient = https.New(config.Username, config.Password)

	return
}
