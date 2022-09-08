package ssh

import (
	"agile-proxy/helper/Go"
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/server/base"
	"encoding/json"
	"fmt"
	gSsh "github.com/gliderlabs/ssh"
	"github.com/pkg/errors"
	sysSsh "golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"time"
)

type ssh struct {
	base.Server
	keyPath string
}

func (s *ssh) Run() (err error) {
	s.Server.Init()
	err = s.listen()
	return
}

func (s *ssh) Close() (err error) {
	if s.Listen != nil {
		err = s.Listen.Close()
	}
	return
}

func (s *ssh) listen() (err error) {
	server := gSsh.Server{
		Addr: net.JoinHostPort(s.Host, s.Port),
		LocalPortForwardingCallback: func(ctx gSsh.Context, destinationHost string, destinationPort uint32) bool {
			return true
		},
	}
	_ = server.SetOption(gSsh.PasswordAuth(s.userInfoAuth()))
	if s.keyPath != "" {
		_ = server.SetOption(gSsh.PublicKeyAuth(s.publicKeyAuth()))
	}
	server.ChannelHandlers = map[string]gSsh.ChannelHandler{
		"direct-tcpip": func(srv *gSsh.Server, conn *sysSsh.ServerConn, newChan sysSsh.NewChannel, ctx gSsh.Context) {
			s.handleDirectRequest(srv, conn, newChan, ctx)
		},
	}

	s.Listen, err = net.Listen("tcp", server.Addr)
	if err != nil {
		err = errors.Wrap(err, "net.Listen")
		return
	}

	errCh := make(chan error)
	Go.Go(func() {
		err := server.Serve(s.Listen)
		if err != nil {
			select {
			case errCh <- err:
			case <-time.After(time.Second * 5):
				log.WarnF("server: %v server.ListenAndServe failed-1: %v", s.Name(), err)
			}
		}
	})

	select {
	case err = <-errCh:
		log.WarnF("server: %v server.ListenAndServe failed-2: %v", s.Name(), err)
	case <-time.After(time.Second * 5):
		log.InfoF("server: %v init successful, listen: %v", s.Name(), server.Addr)
	}

	return
}

func (s *ssh) handleDirectRequest(srv *gSsh.Server, conn *sysSsh.ServerConn, newChan sysSsh.NewChannel, ctx gSsh.Context) {
	d := DirectForward{}
	if err := sysSsh.Unmarshal(newChan.ExtraData(), &d); err != nil {
		log.WarnF("sysSsh.Unmarshal failed: %v", err)
		return
	}

	if srv.LocalPortForwardingCallback == nil || !srv.LocalPortForwardingCallback(ctx, d.DesAddr, d.DesPort) {
		log.Warn("port forwarding is disabled")
		return
	}

	ch, reqs, err := newChan.Accept()
	if err != nil {
		log.WarnF("newChan.Accept failed: %v", err)
		return
	}

	defer ch.Close()
	sshConn := new(Conn)
	sshConn.Channel = ch
	sshConn.localAddr = conn.LocalAddr()
	sshConn.remoteAddr = conn.RemoteAddr()
	Go.Go(func() {
		sysSsh.DiscardRequests(reqs)
	})
	err = s.transport(sshConn, common.StrToBytes(d.DesAddr), common.StrToBytes(fmt.Sprintf("%v", d.DesPort)))
	if err != nil {
		log.WarnF("s.transport %+v", err)
		return
	}
}

func (s *ssh) transport(conn net.Conn, desHost, desPort []byte) (err error) {
	if s.Route != nil {
		err = s.Route.Transport(conn, desHost, desPort)
	} else {
		err = errors.New("Route is nil")
	}
	return
}

func (s *ssh) userInfoAuth() gSsh.PasswordHandler {
	return func(ctx gSsh.Context, password string) bool {
		if ctx.User() != s.Username {
			log.WarnF("server: %v userInfoAuth failed username error", s.Name())
			return false
		}

		if password != s.Password {
			log.WarnF("server: %v userInfoAuth failed password error", s.Name())
			return false
		}

		return true
	}
}

func (s *ssh) publicKeyAuth() gSsh.PublicKeyHandler {
	return func(ctx gSsh.Context, key gSsh.PublicKey) bool {
		data, _ := ioutil.ReadFile(s.keyPath)
		allowed, _, _, _, err := gSsh.ParseAuthorizedKey(data)
		if err != nil {
			log.WarnF("ssh.ParseAuthorizedKey failed: %v", err)
		}
		return gSsh.KeysEqual(key, allowed)
	}
}

func New(jsonConfig json.RawMessage) (obj *ssh, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "new")
		return
	}

	obj = &ssh{
		Server: base.Server{
			Net:        assembly.CreateNet(config.Ip, config.Port, config.Username, config.Password),
			Identity:   assembly.CreateIdentity(config.Name, config.Type),
			Pipeline:   assembly.CreatePipeline(),
			DoneCh:     make(chan struct{}),
			RouteName:  config.RouteName,
			Satellites: config.Satellites,
		},
		keyPath: config.KeyPath,
	}

	return
}
