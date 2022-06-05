package base

import "agile-proxy/modules/plugin"

type Ipc struct {
	OutMsg   plugin.PipelineOutput
	InputMsg plugin.PipelineInput
}

func (i *Ipc) ImplementIpc() {

}
