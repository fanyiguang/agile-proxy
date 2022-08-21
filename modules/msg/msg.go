package msg

import (
	globalConfig "agile-proxy/config"
	"agile-proxy/model"
	"agile-proxy/modules/msg/log"
	sysJson "encoding/json"
	"errors"
	"fmt"
	"strings"

	fileLog "agile-proxy/helper/log"

	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type Msg interface {
	Run() (err error)
	Close() (err error)
	Subscribe(name string, writeCh chan model.ModuleMessage, level string) (chan model.ModuleMessage, string)
	ImplementMsg()
}

func Factory(configs []sysJson.RawMessage) {
	for _, config := range configs {
		var err error
		var message Msg
		switch strings.ToLower(json.Get(config, "type").ToString()) {
		case globalConfig.OutputLog:
			message, err = log.New(config)
		default:
			err = errors.New(fmt.Sprintf("msg type is invalid %v", json.Get(config, "type").ToString()))
		}
		if err != nil {
			fileLog.WarnF("%v", err)
			continue
		}

		msgName := json.Get(config, "name").ToString()
		if msgName != "" {
			if err = message.Run(); err != nil {
				messages[msgName] = message
			} else {
				fileLog.WarnF("msg run failed: %v", err)
			}
		}
	}
}
