package https

import (
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type Client struct {
	username           string
	password           string
	ProxyAuthorization string
	bufferPool         sync.Pool
}

func (c *Client) Handshake(conn net.Conn, target string) (err error) {
	req := &http.Request{
		Method: http.MethodConnect,
		URL:    &url.URL{Host: target},
		Header: make(http.Header),
		Host:   target,
	}

	if c.ProxyAuthorization != "" { // auth
		req.Header.Set("Proxy-Authorization", c.ProxyAuthorization)
	}
	req.Header.Set("Proxy-Connection", "Keep-Alive")
	err = req.Write(conn)
	if err != nil {
		err = errors.Wrap(err, "req.Write")
		return
	}

	buffer := c.bufferPool.Get().([]byte)
	defer c.bufferPool.Put(buffer)
	n, err := conn.Read(buffer)
	if err != nil {
		return err
	}

	log.DebugF("https resp %v", string(buffer[:n]))
	if !strings.Contains(string(buffer[:n]), "200") && !strings.Contains(strings.ToLower(string(buffer[:n])), "connection established") {
		errMsgs := strings.Split(string(buffer[:n]), "\r\n")
		err = errors.New("failed to link to target site. msg:" + errMsgs[0])
	}
	return
}

func New(username, password string) *Client {
	client := &Client{
		username:   username,
		password:   password,
		bufferPool: common.CreateByteBufferSyncPool(1024 * 32),
	}

	if username != "" && password != "" {
		client.ProxyAuthorization = "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", username, password)))
	}
	return client
}
