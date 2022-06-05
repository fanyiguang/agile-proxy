package ssh

import (
	"agile-proxy/config"
	"agile-proxy/helper/Go"
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/modules/client/base"
	"agile-proxy/modules/dialer"
	"agile-proxy/modules/plugin"
	pkgSsh "agile-proxy/pkg/ssh"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"net/url"
	"time"
)

type Ssh struct {
	base.Client
	client           *pkgSsh.Client
	initSuccessfulCh chan struct{}
	initFailedCh     chan struct{}
	doneCh           chan struct{}
	initWorkerCh     chan uint8
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
			err := s.client.HeartBeat()
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
			log.InfoF("server: %v keepAlive end", s.Name())
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
			NetInfo: plugin.NetInfo{
				Host:     _config.Ip,
				Port:     _config.Port,
				Username: _config.Username,
				Password: _config.Password,
			},
			IdentInfo: plugin.IdentInfo{
				ModuleName: _config.Name,
				ModuleType: _config.Type,
			},
			OutMsg: plugin.PipelineOutput{
				Ch: plugin.PipelineOutputCh,
			},
			Mode: _config.Mode,
		},
		rsaPath: _config.RsaPath,
	}

	if _config.DialerName != "" {
		obj.Client.Dialer = dialer.GetDialer(_config.DialerName)
	}
	// 初始化ssh客户端
	obj.client = pkgSsh.New(_config.Ip, _config.Port, pkgSsh.SetUsername(_config.Username), pkgSsh.SetPassword(_config.Password), pkgSsh.SetRsaPath(_config.RsaPath), pkgSsh.SetDialFunc(func(network string, host, port string, timeout time.Duration) (conn net.Conn, err error) {
		return obj.Dialer.DialTimeout(network, host, port, timeout)
	}))

	return
}
