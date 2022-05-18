package socks5

import (
	"bufio"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"nimble-proxy/helper/Go"
	"nimble-proxy/helper/log"
	"nimble-proxy/modules/server"
	"nimble-proxy/modules/transport"
	"nimble-proxy/pkg/socks5"
)

type Socks5 struct {
	server.BaseServer
	Auth int
}

func (s *Socks5) Run() (err error) {
	err = s.listen()
	if err != nil {
		return
	}

	Go.Go(func() {
		s.Accept()
	})

	return
}

func (s *Socks5) Accept() {
	for {
		select {
		case <-s.DoneCh:
			log.InfoF("server: %v accept end", s.Name)
		default:
			conn, err := s.Listen.Accept()
			if err != nil {
				log.WarnF("s.Listen.Accept failed: %v", err)
				continue
			}
			err = s.handler(conn)
			if err != nil {
				log.WarnF("server: %v handler failed: %v", s.Name, err)
			}
		}
	}
}

func (s *Socks5) Transport() {
	//TODO implement me
	panic("implement me")
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
	log.InfoF("server: %v init successful, listen: %v", s.Name, s.Port)
	return
}

func (s *Socks5) handler(conn net.Conn) (err error) {
	defer func() {
		_ = conn.Close()
	}()

	reader := bufio.NewReader(conn)
	_type, _err := reader.ReadByte()
	if _err != nil {
		err = errors.Wrap(_err, "reader.ReadByte")
	}
	socks5.IsSocks5(_type)
	return
	//io.ReadFull(reader)
}

func New(strConfig string) (obj *Socks5, err error) {
	var config Config
	err = json.Unmarshal([]byte(strConfig), &config)
	if err != nil {
		err = errors.Wrap(err, "socks5 new")
		return
	}

	obj = &Socks5{
		BaseServer: server.BaseServer{
			Ip:       config.Ip,
			Port:     config.Port,
			Username: config.Username,
			Password: config.Password,
		},
		Auth: config.Auth,
	}

	if len(config.TransportName) > 0 {
		obj.Transmitter = transport.GetTransport(config.TransportName)
	}

	return
}
