package tool

type MessageList struct {
	Msg chan interface{}
}

func NewMessageList() *MessageList {
	m := &MessageList{
		Msg: make(chan interface{}, 10),
	}
	return m
}
