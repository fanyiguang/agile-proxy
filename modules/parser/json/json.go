package json

import (
	"encoding/json"
	"nimble-proxy/modules/parser"
)

type Json struct {
}

func (j Json) Parser(config []byte) (proxyConfig parser.ProxyConfigInfo, err error) {
	err = json.Unmarshal(config, &proxyConfig)
	return
}
