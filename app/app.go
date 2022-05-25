package app

import (
	"io"
	"nimble-proxy/helper/log"
	"nimble-proxy/modules/client"
	"nimble-proxy/modules/dialer"
	"nimble-proxy/modules/parser"
	"nimble-proxy/modules/server"
	"nimble-proxy/modules/transport"
	"os"
)

func App(configPath string) (err error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return err
	}

	var config []byte
	config, err = io.ReadAll(configFile)
	if err != nil {
		return
	}

	proxyConfig, err := parser.Config(config)
	if err != nil {
		return
	}

	// init log
	log.New(proxyConfig.LogPath, proxyConfig.LogLevel)

	// 固定初始化顺序无法改变
	//依赖关系：server -> transport -> client-> dialer
	dialer.Factory(proxyConfig.DialerConfig)
	client.Factory(proxyConfig.ClientConfig)
	transport.Factory(proxyConfig.TransportConfig)
	servers := server.Factory(proxyConfig.ServerConfig)
	for _, s := range servers {
		err := s.Run()
		if err != nil {
			log.WarnF("%v(%v) run failed: %v", s.Name(), s.Type(), err)
		}
	}

	select {}
}
