package base

import (
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/model"
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/dialer"
	"agile-proxy/modules/satellite"
	"context"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	assembly.Net
	assembly.Identity
	assembly.Pipeline
	model.Satellites
	Dialer       dialer.Dialer
	RoundTripper http.RoundTripper
	DialerName   string
	Mode         int // 0-降级模式（如果有配置连接器且连接器无法使用会走默认网络，默认为降级模式） 1-严格模式（如果有配置连接器且连接器无法使用则直接返回失败）
}

func (s *Client) Dial(network, host, port string) (conn net.Conn, err error) {
	if s.Dialer != nil {
		conn, err = s.Dialer.Dial(network, host, port)
		if err == nil || s.Mode == 1 { // mode=1 严格模式
			return
		}

		if err != nil {
			log.WarnF("s.dialer.Dial failed: %+v", err)
		}
	}

	conn, err = net.Dial(network, net.JoinHostPort(host, port))
	if err != nil {
		err = errors.Wrap(err, "socks5 Dial")
	}
	return
}

func (s *Client) DialTimeout(network, host, port string, timeout time.Duration) (conn net.Conn, err error) {
	if s.Dialer != nil {
		conn, err = s.Dialer.DialTimeout(network, host, port, timeout)
		if err == nil || s.Mode == 1 { // mode=1 严格模式
			return
		}

		if err != nil {
			log.WarnF("s.dialer.Dial failed: %v", err)
		}
	}

	conn, err = net.DialTimeout(network, net.JoinHostPort(host, port), timeout)
	if err != nil {
		err = errors.Wrap(err, "socks5 Dial")
	}
	return
}

func (s *Client) GetStrPort(bPort []byte) string {
	return common.BytesToStr(bPort)
}

func (s *Client) GetRoundTripper() http.RoundTripper {
	return s.RoundTripper
}

func (s *Client) Init() {
	if s.DialerName != "" {
		s.Dialer = dialer.GetDialer(s.DialerName)
	}

	for _, _satellite := range s.Satellites.Satellites {
		_msg := satellite.GetSatellite(_satellite.Name)
		if _msg != nil {
			msgPipeline, level := _msg.Subscribe(s.Name(), s.Pipeline.PipeCh, _satellite.Level)
			s.Subscribe(_satellite.Name, msgPipeline, level)
		} else {
			log.WarnF("%v client get Satellites failed pipeline name: %v", s.Name(), _satellite.Name)
		}
	}
}

func (s *Client) CreateRoundTripper(proxyURL string, dialContext func(ctx context.Context, network, addr string) (net.Conn, error)) (obj *http.Transport, err error) {
	if dialContext == nil {
		d := &net.Dialer{
			Timeout:   12 * time.Second,
			KeepAlive: 30 * time.Second,
		}
		dialContext = d.DialContext
	}
	obj = &http.Transport{
		DialContext:           dialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          200,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if proxyURL != "" {
		_proxy, _err := url.Parse(proxyURL)
		if err != nil {
			err = errors.Wrap(_err, "url.Parse(proxyURL)-1")
			return
		}
		obj.Proxy = http.ProxyURL(_proxy)
	}
	return
}
