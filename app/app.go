package app

import (
	"agile-proxy/config"
	"agile-proxy/helper/Go"
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/helper/process"
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
	"runtime"
	"syscall"
	"time"
)

func App(configPath string, version bool, pprof int) (err error) {
	if version {
		fmt.Printf("agile-proxy %v/%v v%v\n", runtime.GOOS, runtime.GOARCH, config.Version())
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
	// 依赖关系：msg -> server -> route -> client -> dialer
	msg.Factory(proxyConfig.MsgConfig)
	dialer.Factory(proxyConfig.DialerConfig)
	client.Factory(proxyConfig.ClientConfig)
	route.Factory(proxyConfig.RouteConfig)
	server.Factory(proxyConfig.ServerConfig)

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
	case s := <-doneCh:
		log.InfoF("signal done %v", s.String())
	case <-parentProcessDone():
		log.Info("parent done")
	}
}

func parentProcessDone() chan struct{} {
	doneCh := make(chan struct{})
	Go.Go(func() {
		getPpid := os.Getppid()
		ticker := time.NewTicker(time.Second * 15)
		for {
			select {
			case <-ticker.C:
				isRun, err := process.IsRunning(getPpid)
				if err != nil || !isRun {
					log.ErrorF("parent program done: %v %v", isRun, err)
					common.CloseChan(doneCh)
					return
				}
			}
		}
	})
	return doneCh
}

func closeResources() {
	msg.CloseAllMsg()
	server.CloseAllServers()
	route.CloseAllRoutes()
	client.CloseAllClients()
	dialer.CloseAllDialer()
}
