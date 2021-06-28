package tool

type MessageList struct {
	Msg chan string
}

func NewMessageList() *MessageList {
	m := &MessageList{
		Msg: make(chan string),
	}
	return m
}
