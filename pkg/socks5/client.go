package socks5

import (
	"fmt"
	"github.com/pkg/errors"
	"net"
	"nimble-proxy/helper/common"
)

type Client struct {
	username string
	password string
	authMode uint8 // 认证模式 0-允许匿名模式 1-不允许匿名模式
}

type ClientOperation func(client *Client)

func SetClientUsername(username string) ClientOperation {
	return func(client *Client) {
		client.username = username
	}
}

func SetClientPassword(password string) ClientOperation {
	return func(client *Client) {
		client.password = password
	}
}

func SetClientAuth(authMode int) ClientOperation {
	return func(client *Client) {
		client.authMode = uint8(authMode)
	}
}

func NewClient(operates ...ClientOperation) *Client {
	client := new(Client)
	for _, operate := range operates {
		operate(client)
	}
	return client
}

func (c *Client) HandShark(conn net.Conn, desHost, desPort []byte) (err error) {
	usedAuthMode, err := c.handShake(conn)
	if err != nil {
		return
	}

	if usedAuthMode == pass {
		err = c.authentication(conn)
		if err != nil {
			return
		}
	}

	err = c.sendReqInfo(conn, desHost, desPort)
	return
}

func (c *Client) handShake(conn net.Conn) (usedAuthMode uint8, err error) {
	// 目前暂时只支持匿名和密码认证
	switch c.authMode {
	case modeNoAuth:
		_, err = conn.Write(noAuthRequest)
	case modePass:
		_, err = conn.Write(supportPassAuthRequest)
	default:
		err = errors.New("invalid auth_type")
	}
	if err != nil {
		err = errors.Wrap(err, "")
		return
	}

	buffer := make([]byte, 2)
	n, _err := conn.Read(buffer)
	if _err != nil {
		err = errors.Wrap(_err, "c.conn.Read")
		return
	}

	if n < 2 {
		// 可预期错误可以不用打印堆栈，以此来提高性能和日志可读性
		err = errors.New("client.handshake read len < 2")
		return
	}

	if buffer[0] != 5 {
		err = errors.New(fmt.Sprintf("socks5 server response not socks5 protocol data: %#v", buffer))
		return
	}

	// 目前暂时只支持匿名和密码返回的校验
	switch buffer[1] {
	case noAuth:
		usedAuthMode = noAuth
	case pass:
		usedAuthMode = pass
	case errorAuth:
		err = errors.New("socks5 server not support auth mode")
	default:
		err = errors.New("socks5 server response auth mode is invalid")
	}

	return
}

func (c *Client) authentication(conn net.Conn) (err error) {
	reqBuffer := []byte{0x05} // 认证子协商版本（与SOCKS协议版本的0x05无关系，为其他值亦可）
	reqBuffer = append(reqBuffer, byte(len(c.username)))
	reqBuffer = append(reqBuffer, c.username...)
	reqBuffer = append(reqBuffer, byte(len(c.password)))
	reqBuffer = append(reqBuffer, c.password...)
	_, err = conn.Write(reqBuffer)
	if err != nil {
		err = errors.Wrap(err, "c.conn.Write")
		return
	}

	resBuffer := make([]byte, 2)
	n, _err := conn.Read(resBuffer)
	if _err != nil {
		err = errors.Wrap(_err, "c.conn.Read")
		return
	}

	if n < 2 {
		err = errors.New("client.authentication read len < 2")
		return
	}

	if resBuffer[1] != 0 {
		err = errors.New(fmt.Sprintf("client.authentication auth failed. buffur: %#v", resBuffer))
	}
	return
}

func (c *Client) sendReqInfo(conn net.Conn, desHost, desPort []byte) (err error) {
	reqBuffer := []byte{0x05, 0x01, 0x00}
	if ip := net.ParseIP(common.BytesToStr(desHost)); ip != nil {
		reqBuffer = append(reqBuffer, 0x01)
	} else {
		reqBuffer = append(reqBuffer, 0x03)
		reqBuffer = append(reqBuffer, byte(common.GetBytesLen(desHost)))
	}
	reqBuffer = append(reqBuffer, desHost...)
	reqBuffer = append(reqBuffer, desPort...)
	_, err = conn.Write(reqBuffer)

	// 如果响应的type为0x01(ipv4)的长度为：1+1+1+1+4+2
	// 如果响应的type为0x03(ipv6)的长度为：1+1+1+1+16+2 兼容前者
	// 故使用1+1+1+1+16+2
	resBuffer := make([]byte, 1*4+16+2)
	n, _err := conn.Read(resBuffer)
	if _err != nil {
		err = errors.Wrap(_err, "c.conn.Read")
		return
	}

	if n < 2 {
		err = errors.New("client.sendReqInfo read len < 2")
		return
	}

	// 正常需要分析多种情况这边就直接判断成功与否了，偷懒一下。
	if resBuffer[1] != 0 {
		err = errors.New(fmt.Sprintf("client.sendReqInfo send failed. buffur: %#v", resBuffer))
	}
	return
}
