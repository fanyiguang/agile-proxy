package app

import (
	"agile-proxy/helper/Go"
	"agile-proxy/helper/log"
	"agile-proxy/modules/client"
	"agile-proxy/modules/dialer"
	"agile-proxy/modules/ipc"
	"agile-proxy/modules/parser"
	"agile-proxy/modules/server"
	"agile-proxy/modules/transport"
	"io"
	"os"
	"os/signal"
	"syscall"
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
		_s := s
		Go.Go(func() {
			err := _s.Run()
			if err != nil {
				log.WarnF("%v(%v) run failed: %v", _s.Name(), _s.Type(), err)
			}
		})

	}

	wait()
	return
}

func wait() {
	doneCh := make(chan os.Signal, 1)
	signal.Notify(doneCh, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	<-doneCh
	closeResources()
}

func closeResources() {
	server.CloseAllServers()
	transport.CloseAllTransports()
	client.CloseAllClients()
	dialer.CloseAllDialer()
}
