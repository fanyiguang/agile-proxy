package app

import (
	"agile-proxy/config"
	"agile-proxy/helper/Go"
	"agile-proxy/helper/log"
	"agile-proxy/modules/client"
	"agile-proxy/modules/dialer"
	"agile-proxy/modules/msg"
	"agile-proxy/modules/parser"
	"agile-proxy/modules/parser/model"
	"agile-proxy/modules/server"
	"agile-proxy/modules/transport"
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

func App(configPath string, version bool, pprof int) (err error) {
	if version {
		fmt.Printf("agile-proxy v%v\n", config.Version())
		return
	}

	startPprof(pprof)

	proxyConfig, err := parserConfig(configPath)
	if err != nil {
		return err
	}

	// init log
	log.New(proxyConfig.LogPath, proxyConfig.LogLevel)

	// 固定初始化顺序无法改变
	// 依赖关系：server -> transport -> client -> dialer
	dialer.Factory(proxyConfig.DialerConfig)
	client.Factory(proxyConfig.ClientConfig)
	transport.Factory(proxyConfig.TransportConfig)
	_ipc, err := msg.Factory(proxyConfig.MsgConfig)
	if err != nil {
		log.WarnF("msg factory failed: %+v", err)
	} else {
		err = _ipc.Run()
		if err != nil {
			log.WarnF("msg run failed: %+v", err)
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
	closeResources()
	return
}

func startPprof(pprof int) {
	if pprof > 0 {
		Go.Go(func() {
			log.WarnF("start pprof failed: %v", http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", pprof), nil))
		})
	}
}

func parserConfig(configPath string) (proxyConfig model.ProxyConfig, err error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return proxyConfig, err
	}

	var bConfig []byte
	bConfig, err = io.ReadAll(configFile)
	if err != nil {
		return
	}

	proxyConfig, err = parser.Config(bConfig)
	return
}

func wait() {
	doneCh := make(chan os.Signal, 1)
	signal.Notify(doneCh, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	<-doneCh
}

func closeResources() {
	server.CloseAllServers()
	transport.CloseAllTransports()
	client.CloseAllClients()
	dialer.CloseAllDialer()
}
