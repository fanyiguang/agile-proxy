package parser

import (
	"nimble-proxy/modules/parser/json"
	"nimble-proxy/modules/parser/model"
)

var mode = Json

type Parser interface {
	Parser(config []byte) (proxyConfig model.ProxyConfig, err error)
}

func Config(config []byte) (proxyConfig model.ProxyConfig, err error) {
	switch mode {
	case Json:
		parser := new(json.Json)
		return parser.Parser(config)
	default: // default json
		parser := new(json.Json)
		return parser.Parser(config)
	}
}
