package ssh

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"net/url"
	"nimble-proxy/config"
	"nimble-proxy/helper/Go"
	"nimble-proxy/helper/common"
	"nimble-proxy/helper/log"
	"nimble-proxy/modules/client/base"
	"time"
)

type Ssh struct {
	base.Client
	client           *ssh.Client
	config           *ssh.ClientConfig
	initSuccessfulCh chan struct{}
	initFailedCh     chan struct{}
	initWorkerCh     chan uint8
	doneCh           chan struct{}
	rsaPath          string
	network          string
	timeout          int
}

func (s *Ssh) Dial(network string, host, port []byte) (conn net.Conn, err error) {
	err = s.controlCenter()
	if err != nil {
		return
	}

	conn, err = s.client.Dial(network, net.JoinHostPort(string(host), s.GetStrPort(port)))
	if err != nil {
		err = errors.Wrap(err, "s.client.Dial")
	}
	return
}

func (s *Ssh) DialTimeout(network string, host, port []byte, timeout time.Duration) (conn net.Conn, err error) {
	err = s.controlCenter()
	if err != nil {
		return
	}

	resCh := make(chan struct{})
	Go.Go(func() {
		conn, err = s.client.Dial(network, net.JoinHostPort(string(host), s.GetStrPort(port)))
		if err != nil {
			err = errors.Wrap(err, "s.client.Dial")
		}
		common.CloseChan(resCh)
	})

	select {
	case <-resCh:
	case <-time.After(timeout):
		err = dialTimeoutError
	}
	return
}

func (s *Ssh) Close() (err error) {
	common.CloseChan(s.doneCh)
	if s.client != nil {
		_ = s.client.Close()
	}
	return
}

func (s *Ssh) controlCenter() (err error) {
	select {
	case <-s.initWorkerCh:
		err = s.connect()
	case <-s.initSuccessfulCh: // 成功

	case <-s.initFailedCh: // 失败
		err = initFailedError
	case <-time.After(10 * time.Second):
		err = initTimeoutError
	}
	return
}

func (s *Ssh) createConfig() (sshConfig *ssh.ClientConfig, err error) {
	if s.config != nil {
		return s.config, nil
	}

	if s.timeout <= 0 {
		s.timeout = 10
	}
	sshConfig = &ssh.ClientConfig{
		User:            s.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * time.Duration(s.timeout),
	}
	// 优先账户密码认证
	if s.Password != "" {
		sshConfig.Auth = []ssh.AuthMethod{ssh.Password(s.Password)}
		return
	}

	// 其次密钥认证
	if s.rsaPath != "" {
		buffer, _err := ioutil.ReadFile(s.rsaPath)
		if _err != nil {
			err = errors.Wrap(_err, "ioutil.ReadFile")
			return
		}

		signer, _err := ssh.ParsePrivateKey(buffer)
		if _err != nil {
			err = errors.Wrap(_err, "ssh.ParsePrivateKey")
			return
		}

		sshConfig.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
		return
	}

	err = errors.New("password or rsaPath is nil")
	return
}

func (s *Ssh) connect() (err error) {
	s.client, err = s.sshDial()
	if err != nil {
		common.CloseChan(s.initFailedCh) // 快速失败
		return
	}

	Go.Go(func() {
		s.keepAlive()
	})
	return
}

func (s *Ssh) keepAlive() {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			err := s.heartBeat()
			if err != nil {
				log.WarnF("heartBeat failed: %+v", err)
				s.initSuccessfulCh = make(chan struct{})
				err = s.reconnect()
				if err != nil {
					log.WarnF("reconnect failed: %+v", err)
					return
				}

				common.CloseChan(s.initSuccessfulCh)
				return
			}

		case <-s.doneCh:
			log.InfoF("server: %v keepAlive end", s.ClientName)
			return
		}
	}
}

func (s *Ssh) heartBeat() (err error) {
	var conn net.Conn
	for key, _url := range config.GetIpUrls() {
		parse, _err := url.Parse(_url)
		if _err != nil {
			log.WarnF("url: %v url.Parse failed: %v", _url, _err)
			continue
		}
		conn, err = s.client.Dial("tcp", parse.Host)
		if _err == nil { // 正常
			_ = conn.Close()
			return
		}
		if key > 1 { // 三次失败判定为长连接故障
			break
		}
	}
	return
}

func (s *Ssh) sshDial() (client *ssh.Client, err error) {
	sshConfig, err := s.createConfig()
	if err != nil {
		return nil, err
	}

	if s.network == "" {
		s.network = "tcp"
	}
	conn, err := s.Dialer.DialTimeout(s.network, s.Host, s.Port, sshConfig.Timeout)
	if err != nil {
		return nil, err
	}

	c, chans, reqs, _err := ssh.NewClientConn(conn, net.JoinHostPort(s.Host, s.Port), sshConfig)
	if _err != nil {
		err = errors.Wrap(_err, "ssh.NewClientConn")
		return
	}

	client = ssh.NewClient(c, chans, reqs)
	return
}

func (s *Ssh) reconnect() (err error) {
	s.client, err = s.sshDial()
	return
}
