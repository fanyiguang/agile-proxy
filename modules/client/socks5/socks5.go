package socks5

import (
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/client/base"
	"agile-proxy/pkg/socks5"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"time"
)

type Socks5 struct {
	base.Client
	socks5Client *socks5.Client
	authMode     int
}

func (s *Socks5) Dial(network string, host, port []byte) (conn net.Conn, err error) {
	conn, err = s.Client.Dial(network)
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
	conn, err = s.Client.DialTimeout(network, timeout)
	if err != nil {
		return
	}

	err = s.socks5Client.HandShark(conn, host, port)
	if err != nil {
		_ = conn.Close()
	}
	return
}

func (s *Socks5) Run() (err error) {
	s.init()
	return
}

func (s *Socks5) Close() (err error) {
	return
}

func (s *Socks5) init() {
	s.Client.Init()
	s.socks5Client = socks5.NewClient(socks5.SetClientAuth(s.authMode), socks5.SetClientUsername(s.Username), socks5.SetClientPassword(s.Password))
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
