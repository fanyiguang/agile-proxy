package plugin

var (
	OutputCh = make(chan outputMsg)
	InputCh  = make(chan inputMsg)
)

type outputMsg struct {
	Content string
}

type OutputMsg struct {
	OutputMsgCh chan<- outputMsg
}

type inputMsg struct {
	Content string
}

type InputMsg struct {
	InputMsgCh <-chan inputMsg
}
