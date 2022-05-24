package json

import (
	"encoding/json"
	"github.com/pkg/errors"
	"nimble-proxy/modules/parser"
)

type Json struct {
}

func (j Json) Parser(config []byte) (proxyConfig parser.ProxyConfigInfo, err error) {
	_err := json.Unmarshal(config, &proxyConfig)
	if _err != nil {
		err = errors.Wrap(_err, "json.Parser")
	}
	return
}
