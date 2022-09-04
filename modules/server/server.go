package server

import (
	pConfig "agile-proxy/config"
	"agile-proxy/helper/Go"
	"agile-proxy/helper/log"
	"agile-proxy/modules/server/http"
	"agile-proxy/modules/server/https"
	"agile-proxy/modules/server/socks5"
	"agile-proxy/modules/server/ssh"
	"agile-proxy/modules/server/ssl"
	sysJson "encoding/json"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Server interface {
	Run() (err error)
	Name() string
	Type() string
	Close() (err error)
}

func Factory(configs []sysJson.RawMessage) {
	var err error
	var server Server
	for _, config := range configs {
		switch strings.ToLower(json.Get(config, "type").ToString()) {
		case pConfig.Socks5:
			server, err = socks5.New(config)
		case pConfig.Ssl:
			server, err = ssl.New(config)
		case pConfig.Ssh:
			server, err = ssh.New(config)
		case pConfig.Https:
			server, err = https.New(config)
		case pConfig.Http:
			server, err = http.New(config)
		default:
			err = errors.New("type is invalid")
		}
		if err != nil {
			log.WarnF("server init failed: %v", err)
			continue
		}

		servers[server.Name()] = server
		Go.Go(func() {
			err := server.Run()
			if err != nil {
				log.WarnF("%v(%v) run failed: %v", server.Name(), server.Type(), err)
			}
		})
	}
}
