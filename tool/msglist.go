package tool

import "encoding/json"

type MessageList struct {
	msg chan []byte
}

func NewMessageList() *MessageList {
	m := &MessageList{
		msg: make(chan []byte, 10),
	}
	return m
}

func (ml *MessageList) Write(plan float32, txt string) {
	msg_json, _ := json.Marshal(
		struct {
			Plan float32 `json:"plan"`
			Txt  string  `json:"txt"`
		}{Plan: plan, Txt: txt})
	ml.msg <- msg_json
}

func (ml *MessageList) Read() []byte {
	res := <-ml.msg
	return res
}
