package app

import (
	"agile-proxy/config"
	"agile-proxy/helper/Go"
	"agile-proxy/helper/log"
	"agile-proxy/modules/client"
	"agile-proxy/modules/dialer"
	"agile-proxy/modules/msg"
	"agile-proxy/modules/parser"
	"agile-proxy/modules/route"
	"agile-proxy/modules/server"
	"fmt"
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

	proxyConfig, err := parser.Config(configPath)
	if err != nil {
		return err
	}

	// init log
	log.New(proxyConfig.LogPath, proxyConfig.LogLevel)

	// 固定初始化顺序无法改变
	// 依赖关系：msg -> server -> transport -> client -> dialer
	msg.Factory(proxyConfig.MsgConfig)
	dialer.Factory(proxyConfig.DialerConfig)
	client.Factory(proxyConfig.ClientConfig)
	route.Factory(proxyConfig.TransportConfig)
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

func wait() {
	doneCh := make(chan os.Signal, 1)
	signal.Notify(doneCh, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	select {
	case <-doneCh:
		log.InfoF("")
	}
}

func closeResources() {
	server.CloseAllServers()
	route.CloseAllRoutes()
	client.CloseAllClients()
	dialer.CloseAllDialer()
}
