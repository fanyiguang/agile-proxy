package ssh

import (
	"agile-proxy/config"
	"agile-proxy/helper/Go"
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/dialer/base"
	pkgSsh "agile-proxy/proxy/ssh"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"net/url"
	"time"
)

type ssh struct {
	base.Dialer
	assembly.Net
	client           *pkgSsh.Client
	initSuccessfulCh chan struct{}
	initFailedCh     chan struct{}
	doneCh           chan struct{}
	initWorkerCh     chan uint8
	keyPath          string
	network          string
	timeout          int
}

func (s *ssh) Dial(network string, host, port string) (conn net.Conn, err error) {
	err = s.controlCenter()
	if err != nil {
		return
	}

	conn, err = s.client.Dial(network, net.JoinHostPort(host, port))
	if err != nil {
		err = errors.Wrap(err, "s.client.Dial")
	}
	log.DebugF("ssh dialer link status: %v %v", err, net.JoinHostPort(host, port))
	return
}

func (s *ssh) DialTimeout(network string, host, port string, timeout time.Duration) (conn net.Conn, err error) {
	err = s.controlCenter()
	if err != nil {
		return
	}

	resCh := make(chan struct{})
	Go.Go(func() {
		conn, err = s.client.Dial(network, net.JoinHostPort(host, port))
		if err != nil {
			err = errors.Wrap(err, "s.client.Dial")
		}
		common.CloseChan(resCh)
	})

	select {
	case <-resCh:
		log.DebugF("ssh dialer link status: %v %v", err, net.JoinHostPort(host, port))
	case <-time.After(timeout):
		err = dialTimeoutError
	}
	return
}

func (s *ssh) Run() (err error) {
	s.initWorkerCh <- 1
	s.client = pkgSsh.New(s.Host, s.Port, pkgSsh.SetUsername(s.Username), pkgSsh.SetPassword(s.Password), pkgSsh.SetPublicKeyPath(s.keyPath), pkgSsh.SetDialFunc(s.Dialer.BaseDialTimeout))
	return
}

func (s *ssh) Close() (err error) {
	common.CloseChan(s.doneCh)
	if s.client != nil {
		err = s.client.Close()
	}
	return
}

func (s *ssh) controlCenter() (err error) {
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

func (s *ssh) connect() (err error) {
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

func (s *ssh) keepAlive() {
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

func (s *ssh) heartBeat() (err error) {
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

func (s *ssh) reconnect() (err error) {
	err = s.client.Connect()
	return
}

func New(jsonConfig json.RawMessage) (obj *ssh, err error) {
	var _config Config
	err = json.Unmarshal(jsonConfig, &_config)
	if err != nil {
		err = errors.Wrap(err, "new")
		return
	}

	obj = &ssh{
		Dialer: base.Dialer{
			Net:           assembly.CreateNet(_config.Ip, _config.Port, _config.Username, _config.Password),
			Identity:      assembly.CreateIdentity(_config.Name, _config.Type),
			Pipeline:      assembly.CreatePipeline(),
			PipelineInfos: _config.PipelineInfos,
			IFace:         _config.Interface,
		},
		keyPath:          _config.KeyPath,
		initSuccessfulCh: make(chan struct{}),
		initFailedCh:     make(chan struct{}),
		doneCh:           make(chan struct{}),
		initWorkerCh:     make(chan uint8, 1),
	}

	return
}
