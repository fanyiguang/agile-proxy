package http

import (
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/client/base"
	"agile-proxy/proxy/https"
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

func (h *Http) Run() (err error) {
	h.init()
	return
}

func (h *Http) init() {
	h.Client.Init()
	h.httpsClient = https.New(h.Username, h.Password)
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
			Net:           assembly.CreateNet(config.Ip, config.Port, config.Username, config.Password),
			Identity:      assembly.CreateIdentity(config.Name, config.Type),
			Pipeline:      assembly.CreatePipeline(),
			PipelineInfos: config.PipelineInfos,
			Mode:          config.Mode,
			DialerName:    config.DialerName,
		},
	}

	return
}
