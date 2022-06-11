package main

import (
	"agile-proxy/app"
	"flag"
	"log"
)

func main() {
	configPath := flag.String("config", "./config.json", "代理得配置文件路径")
	flag.Parse()
	*configPath = `D:\study\go-objects\my\src\agile-proxy\_example\config_ssh_heartbeat.json`
	err := app.App(*configPath)
	if err != nil {
		log.Printf("app.App failed: %#v", err)
	}
}
