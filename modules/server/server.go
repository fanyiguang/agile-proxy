package server

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"net"
	pConfig "nimble-proxy/config"
	"nimble-proxy/helper/log"
	"nimble-proxy/modules/ipc"
	"nimble-proxy/modules/server/socks5"
	"nimble-proxy/modules/transport"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Server interface {
	Run() (err error)
	Close()
}

type BaseServer struct {
	Name        string
	Type        string
	Ip          string
	Port        string
	Username    string
	Password    string
	DoneCh      chan struct{}
	Listen      net.Listener
	Transmitter transport.Transport
	OutputMsgCh chan<- ipc.Msg
}

func Factory(configs []string) (servers []Server) {
	for _, config := range configs {
		var err error
		var server Server
		switch strings.ToLower(json.Get([]byte(config), "type").ToString()) {
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
