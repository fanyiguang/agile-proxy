package msg

import (
	"agile-proxy/modules/msg/log"
	sysJson "encoding/json"
	jsoniter "github.com/json-iterator/go"
	"strings"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type Msg interface {
	Run() (err error)
	Close() (err error)
	ImplementMsg()
}

func Factory(config sysJson.RawMessage) (obj Msg, err error) {
	switch strings.ToLower(json.Get(config, "type").ToString()) {
	case outputLog:
		fallthrough
	default:
		obj, err = log.New()
	}
	return
}
