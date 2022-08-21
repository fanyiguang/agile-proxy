package log

import (
	"agile-proxy/helper/Go"
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/model"
	"agile-proxy/modules/assembly"
	"agile-proxy/modules/msg/base"
	"encoding/json"

	"github.com/pkg/errors"
)

type outputLog struct {
	base.Msg
	doneCh chan struct{}
}

func (o *outputLog) Run() (err error) {
	Go.Go(func() {
		o.accept()
	})
	return
}

func (o *outputLog) accept() {
	for {
		select {
		case msg := <-o.GetPipeCh():
			log.InfoF("msg log accept module message: %v %v", msg.Name, msg.Content)
		case <-o.doneCh:
			log.InfoF("msg log close")
			return
		}
	}
}

func (o *outputLog) Close() (err error) {
	if o.doneCh != nil {
		common.CloseChan(o.doneCh)
	}
	return
}

func New(jsonConfig json.RawMessage) (obj *outputLog, err error) {
	var config Config
	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		marshalJSON, _ := jsonConfig.MarshalJSON()
		err = errors.Wrap(err, common.BytesToStr(marshalJSON))
		return
	}

	obj = &outputLog{
		Msg: base.Msg{
			Identity: assembly.Identity{
				ModuleName: config.Name,
				ModuleType: config.Type,
			},
			Pipeline: assembly.Pipeline{
				PipeCh:          make(chan model.ModuleMessage),
				SubObjs:         make(map[string]chan model.ModuleMessage),
				RealTimeSubObjs: make(map[string]chan model.ModuleMessage),
			},
		},
		doneCh: make(chan struct{}),
	}
	return
}
