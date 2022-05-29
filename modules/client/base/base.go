package base

import (
	"github.com/pkg/errors"
	"net"
	"nimble-proxy/helper/log"
	"nimble-proxy/modules/dialer"
	"strconv"
	"time"
)

type Client struct {
	Dialer     dialer.Dialer
	Host       string
	Port       string
	Username   string
	Password   string
	ClientName string
	ClientType string
	Mode       int // 0-降级模式（如果有配置连接器且连接器无法使用会走默认网络，默认为降级模式） 1-严格模式（如果有配置连接器且连接器无法使用则直接返回失败）
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
	if len(bPort) < 2 {
		return ""
	}
	return strconv.Itoa(int(bPort[0])<<8 | int(bPort[1]))
}
