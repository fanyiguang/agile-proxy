package base

var (
	OutputCh = make(chan outputMsg)
)

type outputMsg struct {
	Content string
}

type OutputMsg struct {
	OutputMsgCh chan<- outputMsg
}
