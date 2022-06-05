package client

import (
	globalConfig "agile-proxy/config"
	"agile-proxy/helper/log"
	"agile-proxy/modules/client/http"
	"agile-proxy/modules/client/https"
	"agile-proxy/modules/client/socks5"
	"agile-proxy/modules/client/ssh"
	"agile-proxy/modules/client/ssl"
	sysJson "encoding/json"
	"github.com/pkg/errors"
	"net"
	"strings"
	"time"
)

type Client interface {
	Dial(network string, host, port []byte) (conn net.Conn, err error)
	DialTimeout(network string, host, port []byte, timeout time.Duration) (conn net.Conn, err error)
	Close() (err error)
}

func Factory(configs []sysJson.RawMessage) {
	for _, config := range configs {
		var err error
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
			client, err = http.New(config)
		default:
			err = errors.New("type is invalid")
		}
		if err != nil {
			log.WarnF("%#v", err)
			continue
		}

		clientName := json.Get(config, "name").ToString()
		if clientName != "" {
			clients[clientName] = client
		}
	}
	return
}
