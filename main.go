package main

import (
	"agile-proxy/app"
	"flag"
	"log"
)

var (
	configPath = flag.String("config", "./config.json", "代理得配置文件路径")
	version    = flag.Bool("version", false, "版本号")
	pprof      = flag.Int("pprof", 0, "pprof端口号")
)

func main() {
	flag.Parse()
	err := app.App(*configPath, *version, *pprof)
	if err != nil {
		log.Printf("app.App failed: %#v", err)
	}
}
