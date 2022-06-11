package plugin

var (
	PipelineOutputCh = make(chan OutputMsg)
	PipelineInputCh  = make(chan InputMsg)
)

type OutputMsg struct {
	Content    string
	ModuleName string
}

type PipelineOutput struct {
	Ch chan OutputMsg
}

type InputMsg struct {
	Content string
}

type PipelineInput struct {
	Ch chan InputMsg
}
