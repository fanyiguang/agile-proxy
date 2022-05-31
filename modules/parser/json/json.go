package json

import (
	"agile-proxy/modules/parser/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Json struct {
}

func (j Json) Parser(config []byte) (proxyConfig model.ProxyConfig, err error) {
	_err := json.Unmarshal(config, &proxyConfig)
	if _err != nil {
		err = errors.Wrap(_err, "json.Parser")
	}
	return
}
