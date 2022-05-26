package socks5

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"nimble-proxy/helper/Go"
	"nimble-proxy/helper/log"
	"nimble-proxy/modules/ipc"
	"nimble-proxy/modules/server/base"
	"nimble-proxy/modules/transport"
	"nimble-proxy/pkg/socks5"
)

type Socks5 struct {
	base.Server
	AuthMode int
}

func (s *Socks5) Run() (err error) {
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
			log.InfoF("server: %v accept end", s.Name)
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

func (s *Socks5) Close() {
	//TODO implement me
	panic("implement me")
}

func (s *Socks5) listen() (err error) {
	if s.Port == "" {
		err = errors.Wrap(errors.New("server port is nil"), "")
		return
	}

	listen, _err := net.Listen("tcp", net.JoinHostPort(s.Ip, s.Port))
	if _err != nil {
		err = errors.Wrap(_err, "net.Listen")
		return
	}

	s.Listen = listen
	log.InfoF("server: %v init successful, listen: %v", s.Name(), s.Port)
	return
}

func (s *Socks5) handler(conn net.Conn) (err error) {
	defer func() {
		_ = conn.Close()
	}()

	socks5Server := socks5.NewServer(conn, socks5.SetServerAuth(s.AuthMode), socks5.SetServerUsername(s.Username), socks5.SetServerPassword(s.Password))
	err = socks5Server.HandShake()
	if err != nil {
		return
	}

	host, port := socks5Server.GetDesInfo()
	log.DebugF("des host: %v port: %v", string(host), port)
	return s.transport(conn, host, port)
}

func New(jsonConfig json.RawMessage) (obj *Socks5, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "socks5 new")
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
		AuthMode: config.AuthMode,
	}

	if len(config.TransportName) > 0 {
		obj.Transmitter = transport.GetTransport(config.TransportName)
	}

	return
}
