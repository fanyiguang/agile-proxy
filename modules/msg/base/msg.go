package base

import "agile-proxy/modules/assembly"

type Msg struct {
	assembly.Identity
	assembly.Pipeline
	Level int
}

// ImplementMsg 没啥实际意义只是为了区分实现
func (m *Msg) ImplementMsg() {
}
