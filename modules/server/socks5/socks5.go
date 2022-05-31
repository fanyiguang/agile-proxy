package socks5

import (
	"agile-proxy/helper/Go"
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/modules/ipc"
	"agile-proxy/modules/server/base"
	"agile-proxy/modules/transport"
	"agile-proxy/pkg/socks5"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
)

type Socks5 struct {
	base.Server
	socks5Server *socks5.Server
	authMode     int
}

func (s *Socks5) Run() (err error) {
	s.init()
	err = s.listen()
	if err != nil {
		return
	}

	Go.Go(func() {
		s.accept()
	})

	return
}

func (s *Socks5) accept() {
	for {
		select {
		case <-s.DoneCh:
			log.InfoF("server: %v accept end", s.Name())
		default:
			conn, err := s.Listen.Accept()
			if err != nil {
				log.WarnF("s.Listen.accept failed: %v", err)
				continue
			}
			err = s.handler(conn)
			if err != nil {
				log.WarnF("server: %v, handler failed: %+v", s.Name(), err)
			}
		}
	}
}

func (s *Socks5) transport(conn net.Conn, desHost, desPort []byte) (err error) {
	if s.Transmitter != nil {
		err = s.Transmitter.Transport(conn, desHost, desPort)
	} else {
		err = errors.New("Transmitter is nil")
	}
	return
}

func (s *Socks5) Close() (err error) {
	common.CloseChan(s.DoneCh)
	if s.Listen != nil {
		err = s.Listen.Close()
	}
	return
}

func (s *Socks5) listen() (err error) {
	// 可预知的错误，可以通过自定义的错误信息
	// 找到错误位置。所以无需使用wrap。
	if s.Port == "" {
		err = errors.New("server port is nil")
		return
	}

	addr := net.JoinHostPort(s.Ip, s.Port)
	s.Listen, err = net.Listen("tcp", addr)
	if err != nil {
		err = errors.Wrap(err, "net.Listen")
		return
	}

	log.InfoF("server: %v init successful, listen: %v", s.Name(), addr)
	return
}

func (s *Socks5) handler(conn net.Conn) (err error) {
	defer func() {
		_ = conn.Close()
	}()

	host, port, err := s.socks5Server.HandShake(conn)
	if err != nil {
		return
	}

	log.DebugF("des host: %v port: %v", string(host), port)
	return s.transport(conn, host, port)
}

func (s *Socks5) init() {
	s.socks5Server = socks5.NewServer(socks5.SetServerAuth(s.authMode), socks5.SetServerUsername(s.Username), socks5.SetServerPassword(s.Password))
}

func New(jsonConfig json.RawMessage) (obj *Socks5, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "new")
		return
	}

	obj = &Socks5{
		Server: base.Server{
			Ip:          config.Ip,
			Port:        config.Port,
			Username:    config.Username,
			Password:    config.Password,
			ServerName:  config.Name,
			ServerType:  config.Type,
			OutputMsgCh: ipc.OutputCh,
			DoneCh:      make(chan struct{}),
		},
		authMode: config.AuthMode,
	}

	if len(config.TransportName) > 0 {
		obj.Transmitter = transport.GetTransport(config.TransportName)
	}

	return
}
