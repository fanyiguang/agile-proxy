package base

import "agile-proxy/modules/plugin"

type Msg struct {
	OutMsg   plugin.PipelineOutput
	InputMsg plugin.PipelineInput
}

func (m *Msg) ImplementMsg() {
	// 没啥实际意义只是为了区分实现
}
