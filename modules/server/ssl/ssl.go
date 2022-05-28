package ssl

import (
	sysTls "crypto/tls"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"nimble-proxy/helper/Go"
	"nimble-proxy/helper/common"
	"nimble-proxy/helper/log"
	"nimble-proxy/helper/tls"
	"nimble-proxy/modules/ipc"
	"nimble-proxy/modules/server/base"
	"nimble-proxy/modules/transport"
	"nimble-proxy/pkg/socks5"
)

type Ssl struct {
	base.Server
	socks5Server *socks5.Server
	crtPath      string
	keyPath      string
	authMode     int
}

func (s *Ssl) Run() (err error) {
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

func (s *Ssl) Close() (err error) {
	common.CloseChan(s.DoneCh)
	_ = s.Listen.Close()
	return
}

func (s *Ssl) listen() (err error) {
	if s.Port == "" {
		err = errors.New("server port is nil")
		return
	}

	tlsConfig, err := tls.CreateConfig(s.crtPath, s.keyPath)
	if err != nil {
		return
	}

	addr := net.JoinHostPort(s.Ip, s.Port)
	s.Listen, err = sysTls.Listen("tcp", addr, tlsConfig)
	if err != nil {
		err = errors.Wrap(err, "tls.Listen")
		return
	}

	log.InfoF("server: %v init successful, listen: %v", s.Name(), addr)
	return
}

func (s *Ssl) accept() {
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

func (s *Ssl) handler(conn net.Conn) (err error) {
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

func (s *Ssl) transport(conn net.Conn, desHost, desPort []byte) (err error) {
	if s.Transmitter != nil {
		err = s.Transmitter.Transport(conn, desHost, desPort)
	} else {
		err = errors.New("Transmitter is nil")
	}
	return
}

func (s *Ssl) init() {
	s.socks5Server = socks5.NewServer(socks5.SetServerAuth(s.authMode), socks5.SetServerUsername(s.Username), socks5.SetServerPassword(s.Password))
}

func New(jsonConfig json.RawMessage) (obj *Ssl, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "socks5 new")
		return
	}

	obj = &Ssl{
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
		crtPath:  config.CrtPath,
		keyPath:  config.KeyPath,
		authMode: config.AuthMode,
	}

	if len(config.TransportName) > 0 {
		obj.Transmitter = transport.GetTransport(config.TransportName)
	}

	return
}
