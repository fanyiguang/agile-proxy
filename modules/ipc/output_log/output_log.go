package output_log

import (
	"agile-proxy/helper/Go"
	"agile-proxy/helper/common"
	"agile-proxy/helper/log"
	"agile-proxy/modules/ipc/base"
	"agile-proxy/modules/plugin"
)

type OutputLog struct {
	base.Ipc
	doneCh chan struct{}
}

func (o *OutputLog) Run() (err error) {
	Go.Go(func() {
		o.accept()
	})
	return
}

func (o *OutputLog) accept() {
	var msg plugin.OutputMsg
	for {
		select {
		case msg = <-o.OutMsg.Ch:
			log.InfoF("ipc accept module message: %v %v", msg.Module, msg.Content)
		case <-o.doneCh:
			log.InfoF("ipc close")
			return
		}
	}
}

func (o *OutputLog) Close() (err error) {
	if o.doneCh != nil {
		common.CloseChan(o.doneCh)
	}
	return
}

func New() (obj *OutputLog, err error) {
	obj = &OutputLog{
		Ipc: base.Ipc{
			OutMsg: plugin.PipelineOutput{
				Ch: plugin.PipelineOutputCh,
			},
		},
		doneCh: make(chan struct{}),
	}
	return
}
