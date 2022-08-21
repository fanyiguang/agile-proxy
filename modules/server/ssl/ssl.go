package ssl

import (
	"agile-proxy/helper/Go"
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/server/base"
	"agile-proxy/proxy/socks5"
	sysTls "crypto/tls"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
)

type ssl struct {
	base.Server
	assembly.Tls
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
	s.Server.Init()
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
			Net:           assembly.CreateNet(config.Ip, config.Port, config.Username, config.Password),
			Identity:      assembly.CreateIdentity(config.Name, config.Type),
			Pipeline:      assembly.CreatePipeline(),
			DoneCh:        make(chan struct{}),
			TransportName: config.TransportName,
			PipelineInfos: config.PipelineInfos,
		},
		Tls:      assembly.CreateTls(config.CrtPath, config.KeyPath, "", ""),
		authMode: config.AuthMode,
	}

	return
}
