package client

import (
	official "encoding/json"
	"github.com/pkg/errors"
	"net"
	globalConfig "nimble-proxy/config"
	"nimble-proxy/helper/log"
	"nimble-proxy/modules/client/socks5"
	"nimble-proxy/modules/client/ssh"
	"nimble-proxy/modules/client/ssl"
	"strings"
	"time"
)

type Client interface {
	Dial(network string, host, port []byte) (conn net.Conn, err error)
	DialTimeout(network string, host, port []byte, timeout time.Duration) (conn net.Conn, err error)
	Close() (err error)
}

func Factory(configs []official.RawMessage) {
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
