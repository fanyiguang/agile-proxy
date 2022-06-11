package ssh

import (
	"agile-proxy/helper/Go"
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/modules/plugin"
	"agile-proxy/modules/server/base"
	"agile-proxy/modules/transport"
	"encoding/json"
	"fmt"
	"github.com/gliderlabs/ssh"
	"github.com/pkg/errors"
	sysSsh "golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"time"
)

type Ssh struct {
	base.Server
	keyPath string
}

func (s *Ssh) Run() (err error) {
	err = s.listen()
	return
}

func (s *Ssh) Close() (err error) {
	if s.Listen != nil {
		err = s.Listen.Close()
	}
	return
}

func (s *Ssh) listen() (err error) {
	server := ssh.Server{
		Addr: net.JoinHostPort(s.Host, s.Port),
		LocalPortForwardingCallback: func(ctx ssh.Context, destinationHost string, destinationPort uint32) bool {
			return true
		},
	}
	_ = server.SetOption(ssh.PasswordAuth(s.userInfoAuth()))
	if s.keyPath != "" {
		_ = server.SetOption(ssh.PublicKeyAuth(s.publicKeyAuth()))
	}
	server.ChannelHandlers = map[string]ssh.ChannelHandler{
		"direct-tcpip": func(srv *ssh.Server, conn *sysSsh.ServerConn, newChan sysSsh.NewChannel, ctx ssh.Context) {
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

func (s *Ssh) handleDirectRequest(srv *ssh.Server, conn *sysSsh.ServerConn, newChan sysSsh.NewChannel, ctx ssh.Context) {
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

func (s *Ssh) transport(conn net.Conn, desHost, desPort []byte) (err error) {
	if s.Transmitter != nil {
		err = s.Transmitter.Transport(conn, desHost, desPort)
	} else {
		err = errors.New("Transmitter is nil")
	}
	return
}

func (s *Ssh) userInfoAuth() ssh.PasswordHandler {
	return func(ctx ssh.Context, password string) bool {
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

func (s *Ssh) publicKeyAuth() ssh.PublicKeyHandler {
	return func(ctx ssh.Context, key ssh.PublicKey) bool {
		data, _ := ioutil.ReadFile(s.keyPath)
		allowed, _, _, _, err := ssh.ParseAuthorizedKey(data)
		if err != nil {
			log.WarnF("ssh.ParseAuthorizedKey failed: %v", err)
		}
		return ssh.KeysEqual(key, allowed)
	}
}

func New(jsonConfig json.RawMessage) (obj *Ssh, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		err = errors.Wrap(err, "new")
		return
	}

	obj = &Ssh{
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
		keyPath: config.KeyPath,
	}

	if config.TransportName != "" {
		obj.Transmitter = transport.GetTransport(config.TransportName)
	}

	return
}
