package ssh

import (
	"agile-proxy/config"
	"agile-proxy/helper/log"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"net/url"
	"time"
)

type Operation func(*Client)

type Client struct {
	client     *ssh.Client
	config     *ssh.ClientConfig
	dialerFunc func(network string, host, port string, timeout time.Duration) (conn net.Conn, err error)
	host       string
	port       string
	username   string
	password   string
	rsaPath    string
	network    string
	timeout    int
}

func (c *Client) Dial(network string, addr string) (conn net.Conn, err error) {
	return c.client.Dial(network, addr)
}

func (c *Client) Connect() (err error) {
	c.client, err = c.connect()
	return
}

func (c *Client) connect() (client *ssh.Client, err error) {
	sshConfig, err := c.createConfig()
	if err != nil {
		return nil, err
	}

	if c.network == "" {
		c.network = "tcp"
	}

	var conn net.Conn
	if c.dialerFunc != nil {
		conn, err = c.dialerFunc(c.network, c.host, c.port, sshConfig.Timeout)
	} else {
		conn, err = net.DialTimeout(c.network, net.JoinHostPort(c.host, c.port), sshConfig.Timeout)
	}
	if err != nil {
		return nil, err
	}

	sshConn, chans, reqs, _err := ssh.NewClientConn(conn, net.JoinHostPort(c.host, c.port), sshConfig)
	if _err != nil {
		_ = conn.Close()
		_ = errors.Wrap(_err, "ssh.NewClientConn")
		return
	}

	client = ssh.NewClient(sshConn, chans, reqs)
	return
}

func (c *Client) createConfig() (sshConfig *ssh.ClientConfig, err error) {
	if c.config != nil {
		return c.config, nil
	}

	if c.timeout <= 0 {
		c.timeout = 10
	}
	sshConfig = &ssh.ClientConfig{
		User:            c.username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(c.timeout) * time.Second,
	}
	// 优先账户密码认证
	if c.password != "" {
		sshConfig.Auth = []ssh.AuthMethod{ssh.Password(c.password)}
		return
	}

	// 其次密钥认证
	if c.rsaPath != "" {
		buffer, _err := ioutil.ReadFile(c.rsaPath)
		if _err != nil {
			err = errors.Wrap(_err, "ioutil.ReadFile")
			return
		}

		signer, _err := ssh.ParsePrivateKey(buffer)
		if _err != nil {
			err = errors.Wrap(_err, "ssh.ParsePrivateKey")
			return
		}

		sshConfig.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
		return
	}

	err = errors.New("password or rsaPath is nil")
	return
}

func (c *Client) HeartBeat() (err error) {
	var conn net.Conn
	for key, _url := range config.GetIpUrls() {
		parse, _err := url.Parse(_url)
		if _err != nil {
			log.WarnF("url: %v url.Parse failed: %v", _url, _err)
			continue
		}
		conn, err = c.client.Dial("tcp", parse.Host)
		if err == nil { // 正常
			_ = conn.Close()
			return
		}
		if key > 1 { // 三次失败判定为长连接故障
			break
		}
	}
	return
}

func (c *Client) Close() (err error) {
	if c.client != nil {
		err = c.client.Close()
	}
	return
}

func SetUsername(username string) Operation {
	return func(client *Client) {
		client.username = username
	}
}

func SetPassword(password string) Operation {
	return func(client *Client) {
		client.password = password
	}
}

func SetRsaPath(rsaPath string) Operation {
	return func(client *Client) {
		client.rsaPath = rsaPath
	}
}

func SetDialFunc(f func(network string, host, port string, timeout time.Duration) (conn net.Conn, err error)) Operation {
	return func(client *Client) {
		client.dialerFunc = f
	}
}

func SetNetwork(network string) Operation {
	return func(client *Client) {
		client.network = network
	}
}

func SetTimeout(timeout int) Operation {
	return func(client *Client) {
		client.timeout = timeout
	}
}

func New(host, port string, operate ...Operation) *Client {
	obj := &Client{
		host: host,
		port: port,
	}
	for _, op := range operate {
		op(obj)
	}
	return obj
}
