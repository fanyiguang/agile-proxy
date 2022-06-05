package ipc

import (
	"agile-proxy/modules/ipc/output_log"
	sysJson "encoding/json"
	jsoniter "github.com/json-iterator/go"
	"strings"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type Ipc interface {
	ImplementIpc()
	Run() (err error)
	Close() (err error)
}

func Factory(config sysJson.RawMessage) (obj Ipc, err error) {
	switch strings.ToLower(json.Get(config, "type").ToString()) {
	case outputLog:
		fallthrough
	default:
		obj, err = output_log.New()
	}
	return
}
