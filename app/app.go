package app

import (
	"agile-proxy/helper/log"
	"agile-proxy/modules/client"
	"agile-proxy/modules/dialer"
	"agile-proxy/modules/ipc"
	"agile-proxy/modules/parser"
	"agile-proxy/modules/server"
	"agile-proxy/modules/transport"
	"io"
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
	// 依赖关系：server -> transport -> client -> dialer
	dialer.Factory(proxyConfig.DialerConfig)
	client.Factory(proxyConfig.ClientConfig)
	transport.Factory(proxyConfig.TransportConfig)
	_ipc, err := ipc.Factory(proxyConfig.IpcConfig)
	if err != nil {
		log.WarnF("ipc factory failed: %+v", err)
	} else {
		err = _ipc.Run()
		if err != nil {
			log.WarnF("ipc run failed: %+v", err)
		}
	}
	servers := server.Factory(proxyConfig.ServerConfig)
	for _, s := range servers {
		err := s.Run()
		if err != nil {
			log.WarnF("%v(%v) run failed: %v", s.Name(), s.Type(), err)
		}
	}

	select {}
}
