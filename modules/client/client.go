package client

import (
	"github.com/pkg/errors"
	"net"
	globalConfig "nimble-proxy/config"
	"nimble-proxy/helper/log"
	"nimble-proxy/modules/client/socks5"
	"strings"
)

type Client interface {
	Dial(network string, host, port []byte) (conn net.Conn, err error)
	Close()
}

type BaseClient struct {
	Host     string
	Port     string
	Username string
	Password string
	Mode     int // 0-降级模式（如果有配置连接器且连接器无法使用会走默认网络，默认为降级模式） 1-严格模式（如果有配置连接器且连接器无法使用则直接返回失败）
}

func Factory(configs []string) {
	for _, config := range configs {
		var err error
		var client Client
		switch strings.ToLower(json.Get([]byte(config), "type").ToString()) {
		case globalConfig.Socks5:
			client, err = socks5.New(config)
		case globalConfig.Ssl:
		case globalConfig.Ssh:
		default:
			err = errors.New("type is invalid")
		}
		if err != nil {
			log.WarnF("%#v", err)
			continue
		}

		clientName := json.Get([]byte(config), "name").ToString()
		if clientName != "" {
			clients[clientName] = client
		}
	}
	return
}
