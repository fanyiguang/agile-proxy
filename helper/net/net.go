package net

import (
	"agile-proxy/config"
	"agile-proxy/helper/log"
	"agile-proxy/helper/net/ssl"
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

func HttpRequestByProxy(serverType string, serverHost string, serverPort string, username, password string, doRequest func(client *http.Client) (string, error)) (resp string, err error) {
	switch serverType {
	case config.Ssh:
		var client *ssh.Client
		sshConfig := createSshConfig(username, password)
		client, err = ssh.Dial("tcp", net.JoinHostPort(serverHost, serverPort), sshConfig)
		if err != nil {
			err = errors.Wrap(err, "ssh.Dial")
			return
		}
		defer client.Close()

		httpTransport := &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				return client.Dial(network, addr)
			},
		}
		httpClient := &http.Client{Transport: httpTransport}
		defer httpClient.CloseIdleConnections()
		resp, err = doRequest(httpClient)

	case config.Ssl:
		var dialer proxy.Dialer
		sslDial := new(ssl.DialSsl)
		if username == "" || password == "" {
			dialer, err = proxy.SOCKS5("tcp", net.JoinHostPort(serverHost, serverPort), nil, sslDial)
		} else {
			dialer, err = proxy.SOCKS5("tcp", net.JoinHostPort(serverHost, serverPort), &proxy.Auth{User: username, Password: password}, sslDial)
		}

		if err != nil {
			err = errors.Wrap(err, "netProxy.SOCKS5-1")
			return
		}

		httpTransport := &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				return dialer.Dial(network, addr)
			},
		}
		httpClient := &http.Client{Transport: httpTransport, Timeout: 30 * time.Second}
		defer httpClient.CloseIdleConnections()
		resp, err = doRequest(httpClient)

	case config.Socks5:
		var dialer proxy.Dialer
		if username == "" || password == "" {
			dialer, err = proxy.SOCKS5("tcp", net.JoinHostPort(serverHost, serverPort), nil, proxy.Direct)
		} else {
			dialer, err = proxy.SOCKS5("tcp", net.JoinHostPort(serverHost, serverPort), &proxy.Auth{User: username, Password: password}, proxy.Direct)
		}

		if err != nil {
			err = errors.Wrap(err, "netProxy.SOCKS5-2")
			return
		}

		httpTransport := &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				return dialer.Dial(network, addr)
			},
		}
		httpClient := &http.Client{Transport: httpTransport, Timeout: 30 * time.Second}
		defer httpClient.CloseIdleConnections()
		resp, err = doRequest(httpClient)

	case config.Https:
		var _proxy *url.URL
		proxyURL := fmt.Sprintf("http://%s:%s@%s:%s", username, password, serverHost, serverPort)
		_proxy, err = url.Parse(proxyURL)
		if err != nil {
			err = errors.Wrap(err, "url.Parse(proxyURL)-1")
			return
		}

		pool := x509.NewCertPool()
		dialer := &net.Dialer{
			Timeout: 30 * time.Second,
		}
		httpTransport := &http.Transport{
			Proxy: http.ProxyURL(_proxy),
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				return tls.DialWithDialer(dialer, network, addr, &tls.Config{
					RootCAs:            pool,
					InsecureSkipVerify: true,
				})
			},
		}
		httpClient := &http.Client{Transport: httpTransport, Timeout: 30 * time.Second}
		defer httpClient.CloseIdleConnections()
		resp, err = doRequest(httpClient)

	case config.Http:
		var _proxy *url.URL
		proxyURL := fmt.Sprintf("http://%s:%s@%s:%s", username, password, serverHost, serverPort)
		_proxy, err = url.Parse(proxyURL)
		if err != nil {
			err = errors.Wrap(err, "url.Parse(proxyURL)-2")
			return
		}

		httpTransport := &http.Transport{
			Proxy: http.ProxyURL(_proxy),
		}
		httpClient := &http.Client{Transport: httpTransport, Timeout: 30 * time.Second}
		defer httpClient.CloseIdleConnections()
		resp, err = doRequest(httpClient)

	default:
		err = errors.New("this proxy type non-existent")
	}

	if err != nil {
		err = errors.WithMessagef(err, "proxy info: %v %v %v %v %v", serverType, serverHost, serverPort, "username", "password")
	}
	return
}

func HttpRequestByTransport(transport http.RoundTripper, doRequest func(client *http.Client) (string, error)) (resp string, err error) {
	httpClient := &http.Client{Transport: transport, Timeout: 15 * time.Second}
	defer httpClient.CloseIdleConnections()
	resp, err = doRequest(httpClient)
	return
}

func createSshConfig(username, password string) (sshConfig *ssh.ClientConfig) {
	sshConfig = &ssh.ClientConfig{
		User:            username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * time.Duration(30),
	}

	/*选择登录认证类型*/
	sshConfig.Auth = []ssh.AuthMethod{ssh.Password(password)}
	return
}

func CopyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func GetIdlePort() (int, error) {
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "0.0.0.0"))
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		return 0, err
	}

	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

func SplitHostAndPort(host string) (addr, port string, err error) {
	if strings.Contains(host, ":") {
		return net.SplitHostPort(host)
	}

	addr = host
	return
}

func JoinHostAndPort(host, port string) (newHost string) {
	if port != "" {
		newHost = net.JoinHostPort(host, port)
	} else {
		newHost = host
	}
	return
}

func GetExternalIP() (ip string, err error) {
	ipRegExpRule := `(((\d\b|[1-9]\d\b|1\d\d\b|2[0-4]\d|25[0-5]).){3}(\d\b|[1-9]\d\b|1\d\d\b|2[0-4]\d|25[0-5]))`
	checkIpUrls := []string{
		"http://lumtest.com/myip.json",
		"https://checkip.amazonaws.com",
		"http://myipip.net/",
		"https://ipinfo.io/json",
	}
	ipCh := make(chan string, len(checkIpUrls))
	for _, targetIp := range checkIpUrls {
		_targetIp := targetIp
		go func() {
			resp, err := http.Get(_targetIp)
			if err != nil {
				log.WarnF("http.Get failed: %v %v", err, _targetIp)
				return
			}

			if resp.StatusCode != http.StatusOK {
				log.WarnF("http.Get failed status code != 200: %v %v", resp.StatusCode, _targetIp)
				return
			}

			body, _ := ioutil.ReadAll(resp.Body)
			compile := regexp.MustCompile(ipRegExpRule)
			ip := compile.FindString(string(body))
			if ip != "" {
				ipCh <- ip
			}
		}()
	}

	select {
	case ip = <-ipCh:
	case <-time.After(time.Second * 10):
		err = errors.New("timeout")
	}
	return
}

func HttpPost(url string, data []byte, headers map[string]interface{}) (body []byte, err error) {
	var (
		request  *http.Request
		response *http.Response
	)
	request, err = http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return
	}

	for head, value := range headers {
		request.Header.Set(head, value.(string))
	}
	response, err = http.DefaultClient.Do(request)
	if response.Body != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return
	}

	if response.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("status code != 200, %v", response.StatusCode))
		return
	}

	body, err = ioutil.ReadAll(response.Body)
	return
}
