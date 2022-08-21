package json_file

import (
	"agile-proxy/modules/parser/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"io"
	"os"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type JsonFile struct {
}

func (j JsonFile) Parser(config interface{}) (proxyConfig model.ProxyConfig, err error) {
	var ok bool
	var configPath string
	if configPath, ok = config.(string); !ok {
		err = errors.New("config is not string")
		return
	}

	configFile, err := os.Open(configPath)
	if err != nil {
		return proxyConfig, err
	}

	var bConfig []byte
	bConfig, err = io.ReadAll(configFile)
	if err != nil {
		return
	}

	_err := json.Unmarshal(bConfig, &proxyConfig)
	if _err != nil {
		err = errors.Wrap(_err, "json.Parser")
	}
	return
}
