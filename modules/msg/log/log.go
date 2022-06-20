package log

import (
	"agile-proxy/helper/Go"
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/modules/msg/base"
	"agile-proxy/modules/plugin"
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
	var msg plugin.OutputMsg
	for {
		select {
		case msg = <-o.OutMsg.Ch:
			log.InfoF("msg log accept module message: %v %v", msg.ModuleName, msg.Content)
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

func New() (obj *outputLog, err error) {
	obj = &outputLog{
		Msg: base.Msg{
			OutMsg: plugin.PipelineOutput{
				Ch: plugin.PipelineOutputCh,
			},
		},
		doneCh: make(chan struct{}),
	}
	return
}
