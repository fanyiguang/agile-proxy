package model

import "fmt"

type Message struct {
	Action  int
	Content string
}

type ModuleMessage struct {
	Message
	Name string
}

func (m ModuleMessage) String() string {
	return fmt.Sprintf("name: %v action: %v content: %v", m.Name, m.Action, m.Content)
}
