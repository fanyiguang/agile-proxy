package http

import (
	"agile-proxy/helper/log"
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/dialer/base"
	"agile-proxy/pkg/https"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"time"
)

type http struct {
	base.Dialer
	assembly.Net
	httpsClient *https.Client
}

func (h *http) Dial(network string, host, port string) (conn net.Conn, err error) {
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

func (h *http) DialTimeout(network string, host, port string, timeout time.Duration) (conn net.Conn, err error) {
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

func (h *http) Close() (err error) {
	return
}

func (h *http) Run() (err error) {
	err = h.init()
	return
}

func (h *http) init() (err error) {
	h.httpsClient = https.New(h.Username, h.Password)
	return
}

func New(jsonConfig json.RawMessage) (obj *http, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "new")
		return
	}

	obj = &http{
		Dialer: base.Dialer{
			Net:           assembly.CreateNet(config.Ip, config.Port, config.Username, config.Password),
			Identity:      assembly.CreateIdentity(config.Name, config.Type),
			Pipeline:      assembly.CreatePipeline(),
			PipelineInfos: config.PipelineInfos,
			IFace:         config.Interface,
		},
	}

	return
}
