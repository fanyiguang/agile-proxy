package msg

var (
	messages = make(map[string]Msg)
)

func GetMsg(name string) (t Msg) {
	return messages[name]
}

func GetAllMsg() map[string]Msg {
	return messages
}

func CloseAllMsg() {
	for _, message := range messages {
		if message != nil {
			_ = message.Close()
		}
	}
}
