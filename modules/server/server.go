package server

import (
	official "encoding/json"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	pConfig "nimble-proxy/config"
	"nimble-proxy/helper/log"
	"nimble-proxy/modules/server/socks5"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Server interface {
	Run() (err error)
	Name() string
	Type() string
	Close() (err error)
}

func Factory(configs []official.RawMessage) (servers []Server) {
	for _, config := range configs {
		var err error
		var server Server
		switch strings.ToLower(json.Get(config, "type").ToString()) {
		case pConfig.Socks5:
			server, err = socks5.New(config)
		case pConfig.Ssl:
		case pConfig.Ssh:
		default:
			err = errors.New("type is invalid")
		}
		if err != nil {
			log.WarnF("server init failed: %v", err)
			continue
		}

		servers = append(servers, server)
	}
	return
}
