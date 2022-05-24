package parser

import "nimble-proxy/modules/parser/json"

var mode = Json

type Parser interface {
	Parser(config []byte) (proxyConfig ProxyConfigInfo, err error)
}

func Config(config []byte) (proxyConfig ProxyConfigInfo, err error) {
	switch mode {
	case Json:
		parser := new(json.Json)
		return parser.Parser(config)
	default: // default json
		parser := new(json.Json)
		return parser.Parser(config)
	}
}
