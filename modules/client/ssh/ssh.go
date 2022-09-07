package ssh

import (
	"agile-proxy/helper/Go"
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/client/base"
	pkgSsh "agile-proxy/proxy/ssh"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"time"
)

type Ssh struct {
	base.Client
	client           *pkgSsh.Client
	initSuccessfulCh chan struct{}
	initFailedCh     chan struct{}
	doneCh           chan struct{}
	initWorkerCh     chan uint8
	keyPath          string
	network          string
	timeout          int
}

func (s *Ssh) Dial(network string, host, port []byte) (conn net.Conn, err error) {
	err = s.controlCenter()
	if err != nil {
		return
	}

	conn, err = s.client.Dial(network, net.JoinHostPort(common.BytesToStr(host), s.GetStrPort(port)))
	if err != nil {
		err = errors.Wrap(err, net.JoinHostPort(common.BytesToStr(host), s.GetStrPort(port)))
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
		conn, err = s.client.Dial(network, net.JoinHostPort(common.BytesToStr(host), s.GetStrPort(port)))
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

func (s *Ssh) Run() (err error) {
	s.Client.Init()
	err = s.createRoundTripper()

	// 添加初始化的令牌，只有一个请求会进入SSH初始化逻辑
	s.initWorkerCh <- 1
	// 初始化ssh客户端
	s.client = pkgSsh.New(s.Host, s.Port, pkgSsh.SetUsername(s.Username), pkgSsh.SetPassword(s.Password), pkgSsh.SetPublicKeyPath(s.keyPath), pkgSsh.SetDialFunc(s.Client.DialTimeout))
	return
}

func (s *Ssh) Close() (err error) {
	common.CloseChan(s.doneCh)
	if s.client != nil {
		err = s.client.Close()
	}
	return
}

func (s *Ssh) controlCenter() (err error) {
	select {
	case <-s.initWorkerCh: // 只会有一个请求进去初始化
		err = s.connect()
	case <-s.initSuccessfulCh: // 成功

	case <-s.initFailedCh: // 失败
		err = initFailedError
	case <-time.After(10 * time.Second):
		err = initTimeoutError
	}
	return
}

func (s *Ssh) connect() (err error) {
	err = s.client.Connect()
	if err != nil {
		// 初始化ssh失败，重新添加初始化的令牌，
		common.ReliableChanSend(s.initWorkerCh, 1)
		return
	}

	// 初始化成功打开控制阀门
	common.CloseChan(s.initSuccessfulCh)
	Go.Go(func() {
		s.keepAlive()
	})
	return
}

func (s *Ssh) createRoundTripper() (err error) {
	s.RoundTripper, err = s.CreateRoundTripper("", func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
		deadline, ok := ctx.Deadline()
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return conn, err
		}

		if ok {
			now := time.Now()
			if deadline.After(now) {
				conn, err = s.DialTimeout(network, common.StrToBytes(host), common.StrToBytes(port), deadline.Sub(now))
			} else {
				err = http.ErrHandlerTimeout
			}
		} else {
			conn, err = s.Dial(network, common.StrToBytes(host), common.StrToBytes(port))
		}
		return
	})
	return
}

func (s *Ssh) keepAlive() {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			err := s.client.HeartBeat()
			if err != nil {
				log.WarnF("heartBeat failed: %+v", err)
				s.initSuccessfulCh = make(chan struct{})
				err = s.reconnect()
				if err != nil {
					log.WarnF("reconnect failed: %+v", err)
					continue
				}

				common.CloseChan(s.initSuccessfulCh)
			}

		case <-s.doneCh:
			log.InfoF("server: %v keepAlive end", s.Name())
			return
		}
	}
}

func (s *Ssh) reconnect() (err error) {
	err = s.client.Connect()
	return
}

func New(jsonConfig json.RawMessage) (obj *Ssh, err error) {
	var _config Config
	err = json.Unmarshal(jsonConfig, &_config)
	if err != nil {
		err = errors.Wrap(err, "new")
		return
	}

	obj = &Ssh{
		Client: base.Client{
			Net:        assembly.CreateNet(_config.Ip, _config.Port, _config.Username, _config.Password),
			Identity:   assembly.CreateIdentity(_config.Name, _config.Type),
			Pipeline:   assembly.CreatePipeline(),
			Satellites: _config.Satellites,
			Mode:       _config.Mode,
			DialerName: _config.DialerName,
		},
		keyPath:          _config.KeyPath,
		initSuccessfulCh: make(chan struct{}),
		initFailedCh:     make(chan struct{}),
		doneCh:           make(chan struct{}),
		initWorkerCh:     make(chan uint8, 1),
	}

	return
}
