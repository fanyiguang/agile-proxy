package socks5

import (
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/dialer/base"
	pkgSocks5 "agile-proxy/proxy/socks5"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"time"
)

type socks5 struct {
	base.Dialer
	socks5Client *pkgSocks5.Client
	authMode     int
}

func (s *socks5) Dial(network string, host, port string) (conn net.Conn, err error) {
	conn, err = s.BaseDial(network, s.Host, s.Port)
	if err != nil {
		return
	}

	err = s.socks5Client.HandShark(conn, common.StrToBytes(host), common.StrToBytes(port))
	if err != nil {
		_ = conn.Close()
	}
	log.DebugF("socks5 dialer link status: %v %v", err, net.JoinHostPort(host, port))
	return
}

func (s *socks5) DialTimeout(network string, host, port string, timeout time.Duration) (conn net.Conn, err error) {
	conn, err = s.BaseDialTimeout(network, s.Host, s.Port, timeout)
	if err != nil {
		return
	}

	err = s.socks5Client.HandShark(conn, common.StrToBytes(host), common.StrToBytes(port))
	if err != nil {
		_ = conn.Close()
	}
	return
}

func (s *socks5) Run() (err error) {
	s.socks5Client = pkgSocks5.NewClient(pkgSocks5.SetClientAuth(s.authMode), pkgSocks5.SetClientUsername(s.Username), pkgSocks5.SetClientPassword(s.Password))
	return
}

func (s *socks5) Close() (err error) {
	return
}

func New(jsonConfig json.RawMessage) (obj *socks5, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "new")
		return
	}

	obj = &socks5{
		Dialer: base.Dialer{
			Net:        assembly.CreateNet(config.Ip, config.Port, config.Username, config.Password),
			Identity:   assembly.CreateIdentity(config.Name, config.Type),
			Pipeline:   assembly.CreatePipeline(),
			Satellites: config.Satellites,
			IFace:      config.Interface,
		},
		authMode: config.AuthMode,
	}

	return
}
