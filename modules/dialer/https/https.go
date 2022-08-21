package https

import (
	"agile-proxy/helper/log"
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/dialer/base"
	pkgHttps "agile-proxy/pkg/https"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"time"
)

type https struct {
	base.Dialer
	assembly.Tls
	httpsClient *pkgHttps.Client
}

func (h *https) Dial(network string, host, port string) (conn net.Conn, err error) {
	conn, err = h.BaseDial(network, h.Host, h.Port)
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

	err = h.httpsClient.Handshake(conn, net.JoinHostPort(host, port))
	if err != nil {
		_ = conn.Close()
	}
	log.DebugF("https dialer link status: %v %v", err, net.JoinHostPort(host, port))
	return
}

func (h *https) DialTimeout(network string, host, port string, timeout time.Duration) (conn net.Conn, err error) {
	conn, err = h.BaseDialTimeout(network, h.Host, h.Port, timeout)
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

	err = h.httpsClient.Handshake(conn, net.JoinHostPort(host, port))
	if err != nil {
		_ = conn.Close()
	}
	log.DebugF("https dialer link status: %v %v", err, net.JoinHostPort(host, port))
	return
}

func (h *https) Run() (err error) {
	err = h.init()
	return
}

func (h *https) Close() (err error) {
	return
}

func (h *https) init() (err error) {
	h.httpsClient = pkgHttps.New(h.Username, h.Password)
	return
}

func New(jsonConfig json.RawMessage) (obj *https, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "new")
		return
	}

	obj = &https{
		Dialer: base.Dialer{
			Net:           assembly.CreateNet(config.Ip, config.Port, config.Username, config.Password),
			Identity:      assembly.CreateIdentity(config.Name, config.Type),
			Pipeline:      assembly.CreatePipeline(),
			PipelineInfos: config.PipelineInfos,
			IFace:         config.Interface,
		},
		Tls: assembly.CreateTls(config.CrtPath, config.KeyPath, config.CaPath, config.ServerName),
	}

	return
}
