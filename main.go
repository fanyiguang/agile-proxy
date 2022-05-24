package main

import (
	"flag"
	"log"
	"nimble-proxy/app"
)

func main() {
	configPath := flag.String("config", "./config.json", "代理得配置文件路径")
	flag.Parse()
	*configPath = `D:\study\go-objects\my\src\nimble-proxy\_example\config-test.json`
	err := app.App(*configPath)
	if err != nil {
		log.Printf("app.App failed: %#v", err)
	}
}
