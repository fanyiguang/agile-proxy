package client

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	globalConfig "nimble-proxy/config"
	"nimble-proxy/helper/log"
	"nimble-proxy/modules/client/socks5"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Client interface {
	Dial()
	Close()
}

type BaseClient struct {
	Ip       string
	Port     string
	Username string
	Password string
}

func Factory(configs []string) (clients []Client) {
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
			log.WarnF("client init failed: %v", err)
			continue
		}

		clients = append(clients, client)
	}
	return
}
