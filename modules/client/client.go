package client

import (
	globalConfig "agile-proxy/config"
	"agile-proxy/helper/log"
	"agile-proxy/modules/client/direct"
	"agile-proxy/modules/client/http"
	"agile-proxy/modules/client/https"
	"agile-proxy/modules/client/socks5"
	"agile-proxy/modules/client/ssh"
	"agile-proxy/modules/client/ssl"
	sysJson "encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net"
	sysHttp "net/http"
	"strings"
	"time"
)

type Client interface {
	Run() (err error)
	Name() string
	Dial(network string, host, port []byte) (conn net.Conn, err error)
	DialTimeout(network string, host, port []byte, timeout time.Duration) (conn net.Conn, err error)
	GetRoundTripper() sysHttp.RoundTripper
	Close() (err error)
}

func Factory(configs []sysJson.RawMessage) {
	var err error
	var clientName string
	for _, config := range configs {
		var client Client
		switch strings.ToLower(json.Get(config, "type").ToString()) {
		case globalConfig.Socks5:
			client, err = socks5.New(config)
		case globalConfig.Ssl:
			client, err = ssl.New(config)
		case globalConfig.Ssh:
			client, err = ssh.New(config)
		case globalConfig.Https:
			client, err = https.New(config)
		case globalConfig.Http:
			client, err = http.New(config)
		case globalConfig.Direct:
			client, err = direct.New(config)
		default:
			err = errors.New(fmt.Sprintf("type is invalid %v", strings.ToLower(json.Get(config, "type").ToString())))
		}
		if err != nil {
			log.WarnF("%#v", err)
			continue
		}

		clientName = json.Get(config, "name").ToString()
		if err = client.Run(); err == nil {
			clients[clientName] = client
		} else {
			log.WarnF("%v client run failed: %v", clientName, err)
		}
	}
	return
}
