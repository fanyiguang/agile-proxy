package client

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"net"
	globalConfig "nimble-proxy/config"
	"nimble-proxy/helper/log"
	"nimble-proxy/modules/client/socks5"
	"strings"
)

var (
	json    = jsoniter.ConfigCompatibleWithStandardLibrary
	clients = make(map[string]Client)
)

type Client interface {
	Dial(host, port []byte) (conn net.Conn, err error)
	Close()
}

type BaseClient struct {
	Ip       string
	Port     string
	Username string
	Password string
}

func GetClient(name string) (t Client) {
	return clients[name]
}

func GetAllClients() map[string]Client {
	return clients
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
