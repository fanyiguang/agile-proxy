package socks5

import (
	"agile-proxy/helper/common"
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/client/base"
	"agile-proxy/proxy/socks5"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"time"
)

type Socks5 struct {
	base.Client
	socks5Client *socks5.Client
	authMode     int
}

func (s *Socks5) Dial(network string, host, port []byte) (conn net.Conn, err error) {
	conn, err = s.Client.Dial(network, s.Host, s.Port)
	if err != nil {
		return
	}

	err = s.socks5Client.HandShark(conn, host, port)
	if err != nil {
		_ = conn.Close()
	}
	return
}

func (s *Socks5) DialTimeout(network string, host, port []byte, timeout time.Duration) (conn net.Conn, err error) {
	conn, err = s.Client.DialTimeout(network, s.Host, s.Port, timeout)
	if err != nil {
		return
	}

	err = s.socks5Client.HandShark(conn, host, port)
	if err != nil {
		_ = conn.Close()
	}
	return
}

func (s *Socks5) createRoundTripper() (err error) {
	s.RoundTripper, err = s.CreateRoundTripper("", func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
		deadline, ok := ctx.Deadline()
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return conn, err
		}

		if ok {
			now := time.Now()
			if deadline.After(now) {
				conn, err = s.DialTimeout(network, common.StrToBytes(host), common.StrToBytes(port), deadline.Sub(now))
			} else {
				err = http.ErrHandlerTimeout
			}
		} else {
			conn, err = s.Dial(network, common.StrToBytes(host), common.StrToBytes(port))
		}
		return
	})
	return
}

func (s *Socks5) Run() (err error) {
	s.Client.Init()
	s.socks5Client = socks5.NewClient(socks5.SetClientAuth(s.authMode), socks5.SetClientUsername(s.Username), socks5.SetClientPassword(s.Password))
	err = s.createRoundTripper()
	return
}

func (s *Socks5) Close() (err error) {
	return
}

func New(jsonConfig json.RawMessage) (obj *Socks5, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "new")
		return
	}

	obj = &Socks5{
		Client: base.Client{
			Net:           assembly.CreateNet(config.Ip, config.Port, config.Username, config.Password),
			Identity:      assembly.CreateIdentity(config.Name, config.Type),
			Pipeline:      assembly.CreatePipeline(),
			PipelineInfos: config.PipelineInfos,
			Mode:          config.Mode,
			DialerName:    config.DialerName,
		},
		authMode: config.AuthMode,
	}

	return
}
