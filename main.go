package main

import (
	"log"
	"nimble-proxy/app"
)

func main() {
	err := app.App()
	if err != nil {
		log.Printf("app.App failed: %#v", err)
	}
}
