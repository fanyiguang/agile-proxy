package https

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"net/url"
)

type Client struct {
	username           string
	password           string
	ProxyAuthorization string
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

	resp, _err := http.ReadResponse(bufio.NewReader(conn), req)
	if _err != nil {
		err = errors.Wrap(_err, "req.ReadResponse")
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("proxy server response code not 200: %v, %v", resp.StatusCode, resp.Status))
	}
	return
}

func New(username, password string) *Client {
	client := &Client{
		username: username,
		password: password,
	}

	if username != "" && password != "" {
		client.ProxyAuthorization = "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", username, password)))
	}
	return client
}
