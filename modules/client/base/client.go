package base

import (
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/modules/dialer"
	"agile-proxy/modules/plugin"
	"github.com/pkg/errors"
	"net"
	"time"
)

type Client struct {
	plugin.Net
	plugin.Identity
	OutMsg plugin.PipelineOutput
	Dialer dialer.Dialer
	Mode   int // 0-降级模式（如果有配置连接器且连接器无法使用会走默认网络，默认为降级模式） 1-严格模式（如果有配置连接器且连接器无法使用则直接返回失败）
}

func (s *Client) Dial(network string) (conn net.Conn, err error) {
	if s.Dialer != nil {
		conn, err = s.Dialer.Dial(network, s.Host, s.Port)
		if err == nil || s.Mode == 1 { // mode=1 严格模式
			return
		}

		if err != nil {
			log.WarnF("s.dialer.Dial failed: %v", err)
		}
	}

	conn, err = net.Dial(network, net.JoinHostPort(s.Host, s.Port))
	if err != nil {
		err = errors.Wrap(err, "socks5 Dial")
	}
	return
}

func (s *Client) DialTimeout(network string, timeout time.Duration) (conn net.Conn, err error) {
	if s.Dialer != nil {
		conn, err = s.Dialer.DialTimeout(network, s.Host, s.Port, timeout)
		if err == nil || s.Mode == 1 { // mode=1 严格模式
			return
		}

		if err != nil {
			log.WarnF("s.dialer.Dial failed: %v", err)
		}
	}

	conn, err = net.Dial(network, net.JoinHostPort(s.Host, s.Port))
	if err != nil {
		err = errors.Wrap(err, "socks5 Dial")
	}
	return
}

func (s *Client) GetStrPort(bPort []byte) string {
	return common.BytesToStr(bPort)
}
