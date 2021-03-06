package ssl

import (
	"agile-proxy/helper/Go"
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/modules/plugin"
	"agile-proxy/modules/server/base"
	"agile-proxy/modules/transport"
	"agile-proxy/pkg/socks5"
	sysTls "crypto/tls"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
)

type ssl struct {
	base.Server
	plugin.Tls
	socks5Server *socks5.Server
	authMode     int
}

func (s *ssl) Run() (err error) {
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

func (s *ssl) Close() (err error) {
	common.CloseChan(s.DoneCh)
	if s.Listen != nil {
		err = s.Listen.Close()
	}
	return
}

func (s *ssl) listen() (err error) {
	if s.Port == "" {
		err = errors.New("server port is nil")
		return
	}

	tlsConfig, err := s.CreateServerTlsConfig()
	if err != nil {
		return
	}

	addr := net.JoinHostPort(s.Host, s.Port)
	s.Listen, err = sysTls.Listen("tcp", addr, tlsConfig)
	if err != nil {
		err = errors.Wrap(err, "tls.Listen")
		return
	}

	log.InfoF("server: %v init successful, listen: %v", s.Name(), addr)
	return
}

func (s *ssl) accept() {
	for {
		select {
		case <-s.DoneCh:
			log.InfoF("server: %v accept end", s.Name())
			return
		default:
			conn, err := s.Listen.Accept()
			if err != nil {
				log.WarnF("s.Listen.accept failed: %v", err)
				continue
			}
			Go.Go(func() {
				err = s.handler(conn)
				if err != nil {
					log.WarnF("server: %v, handler failed: %+v", s.Name(), err)
				}
			})
		}
	}
}

func (s *ssl) handler(conn net.Conn) (err error) {
	defer func() {
		_ = conn.Close()
	}()

	host, port, err := s.socks5Server.HandShake(conn)
	if err != nil {
		return
	}

	//log.DebugF("des host: %v port: %v", string(host), string(port))
	return s.transport(conn, host, port)
}

func (s *ssl) transport(conn net.Conn, desHost, desPort []byte) (err error) {
	if s.Transmitter != nil {
		err = s.Transmitter.Transport(conn, desHost, desPort)
	} else {
		err = errors.New("Transmitter is nil")
	}
	return
}

func (s *ssl) init() {
	s.socks5Server = socks5.NewServer(socks5.SetServerAuth(s.authMode), socks5.SetServerUsername(s.Username), socks5.SetServerPassword(s.Password))
}

func New(jsonConfig json.RawMessage) (obj *ssl, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "new")
		return
	}

	obj = &ssl{
		Server: base.Server{
			Net: plugin.Net{
				Host:     config.Ip,
				Port:     config.Port,
				Username: config.Username,
				Password: config.Password,
			},
			Identity: plugin.Identity{
				ModuleName: config.Name,
				ModuleType: config.Type,
			},
			OutMsg: plugin.PipelineOutput{
				Ch: plugin.PipelineOutputCh,
			},
			DoneCh: make(chan struct{}),
		},
		Tls: plugin.Tls{
			CrtPath: config.CrtPath,
			KeyPath: config.KeyPath,
		},
		authMode: config.AuthMode,
	}

	if len(config.TransportName) > 0 {
		obj.Transmitter = transport.GetTransport(config.TransportName)
	}

	return
}
