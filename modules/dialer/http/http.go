package http

import (
	"agile-proxy/helper/log"
	"agile-proxy/modules/dialer/base"
	"agile-proxy/modules/plugin"
	"agile-proxy/pkg/https"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"time"
)

type Http struct {
	base.Dialer
	plugin.Net
	httpsClient *https.Client
}

func (h *Http) Dial(network string, host, port string) (conn net.Conn, err error) {
	conn, err = h.BaseDial(network, h.Host, h.Port)
	if err != nil {
		return
	}

	err = h.httpsClient.Handshake(conn, net.JoinHostPort(host, port))
	if err != nil {
		_ = conn.Close()
	}
	log.DebugF("http dialer link status: %v %v", err, net.JoinHostPort(host, port))
	return
}

func (h *Http) DialTimeout(network string, host, port string, timeout time.Duration) (conn net.Conn, err error) {
	conn, err = h.BaseDialTimeout(network, h.Host, h.Port, timeout)
	if err != nil {
		return
	}

	err = h.httpsClient.Handshake(conn, net.JoinHostPort(host, port))
	if err != nil {
		_ = conn.Close()
	}
	log.DebugF("http dialer link status: %v %v", err, net.JoinHostPort(host, port))
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
		Dialer: base.Dialer{
			Identity: plugin.Identity{
				ModuleName: config.Name,
				ModuleType: config.Type,
			},
			OutMsg: plugin.PipelineOutput{
				Ch: plugin.PipelineOutputCh,
			},
		},
		Net: plugin.Net{
			Host:     config.Ip,
			Port:     config.Port,
			Username: config.Username,
			Password: config.Password,
		},
	}

	obj.httpsClient = https.New(config.Username, config.Password)
	return
}
