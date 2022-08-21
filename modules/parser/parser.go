package parser

import (
	"agile-proxy/modules/parser/json_file"
	"agile-proxy/modules/parser/model"
)

var mode = JsonFile

type Parser interface {
	Parser(config []byte) (proxyConfig model.ProxyConfig, err error)
}

func Config(config interface{}) (proxyConfig model.ProxyConfig, err error) {
	switch mode {
	case JsonFile:
		fallthrough
	default: // default jsonFile
		parser := new(json_file.JsonFile)
		return parser.Parser(config)
	}
}
