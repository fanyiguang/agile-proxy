package output_log

import (
	"agile-proxy/helper/Go"
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/modules/ipc/base"
	"agile-proxy/modules/plugin"
)

type outputLog struct {
	base.Ipc
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
			log.InfoF("ipc accept module message: %v %v", msg.ModuleName, msg.Content)
		case <-o.doneCh:
			log.InfoF("ipc close")
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
		Ipc: base.Ipc{
			OutMsg: plugin.PipelineOutput{
				Ch: plugin.PipelineOutputCh,
			},
		},
		doneCh: make(chan struct{}),
	}
	return
}
