package https

import (
	"agile-proxy/helper/common"
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

type Https struct {
	base.Client
	assembly.Tls
	httpsClient *https.Client
}

func (h *Https) Dial(network string, host, port []byte) (conn net.Conn, err error) {
	conn, err = h.Client.Dial(network, h.Host, h.Port)
	if err != nil {
		return
	}

	config, err := h.CreateClientTlsConfig()
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
	conn, err = h.Client.DialTimeout(network, h.Host, h.Port, timeout)
	if err != nil {
		return
	}

	config, err := h.CreateClientTlsConfig()
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

func (h *Https) createRoundTripper() (err error) {
	proxyURL := fmt.Sprintf("http://%s:%s@%s:%s", h.Username, h.Password, h.Host, h.Port)
	config, err := h.CreateClientTlsConfig()
	if err != nil {
		return err
	}

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

		conn, err = h.Handshake(ctx, conn, config)
		return
	})
	return
}

func (h *Https) Close() (err error) {
	return
}

func (h *Https) Run() (err error) {
	h.Client.Init()
	h.httpsClient = https.New(h.Username, h.Password)
	err = h.createRoundTripper()
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
			Net:           assembly.CreateNet(config.Ip, config.Port, config.Username, config.Password),
			Identity:      assembly.CreateIdentity(config.Name, config.Type),
			Pipeline:      assembly.CreatePipeline(),
			PipelineInfos: config.PipelineInfos,
			Mode:          config.Mode,
			DialerName:    config.DialerName,
		},
		Tls: assembly.CreateTls(config.CrtPath, config.KeyPath, config.CaPath, config.ServerName),
	}

	return
}
