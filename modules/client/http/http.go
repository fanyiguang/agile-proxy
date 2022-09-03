package http

import (
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/client/base"
	"agile-proxy/proxy/https"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"time"
)

type Http struct {
	base.Client
	httpsClient *https.Client
}

func (h *Http) Dial(network string, host, port []byte) (conn net.Conn, err error) {
	conn, err = h.Client.Dial(network, h.Host, h.Port)
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
	conn, err = h.Client.DialTimeout(network, h.Host, h.Port, timeout)
	if err != nil {
		return
	}

	err = h.httpsClient.Handshake(conn, net.JoinHostPort(string(host), h.GetStrPort(port)))
	if err != nil {
		_ = conn.Close()
	}
	return
}

func (h *Http) createRoundTripper() (err error) {
	proxyURL := fmt.Sprintf("http://%s:%s@%s:%s", h.Username, h.Password, h.Host, h.Port)
	h.RoundTripper, err = h.CreateRoundTripper(proxyURL, func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
		deadline, ok := ctx.Deadline()
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return conn, err
		}

		if ok {
			now := time.Now()
			if deadline.After(now) {
				conn, err = h.Client.DialTimeout(network, host, port, deadline.Sub(now))
			} else {
				err = http.ErrHandlerTimeout
			}
		} else {
			conn, err = h.Client.Dial(network, host, port)
		}
		return
	})
	return
}

func (h *Http) Close() (err error) {
	return
}

func (h *Http) Run() (err error) {
	h.Client.Init()
	h.httpsClient = https.New(h.Username, h.Password)
	err = h.createRoundTripper()
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
