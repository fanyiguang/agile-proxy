package base

var (
	InputCh = make(chan inputMsg)
)

type inputMsg struct {
	Content string
}

type InputMsg struct {
	InputMsgCh <-chan inputMsg
}
