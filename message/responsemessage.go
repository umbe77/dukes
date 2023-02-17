// Copyright (c) 2023 Robeto Ughi
// 
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package message

type Status uint8

const (
	BeginResp Status = iota
	EndResp
	OK
	BadFormat
	Error
	Nope
)

type ResponseMessage struct {
	St     Status
	Params []MessageParam
}

func (m ResponseMessage) ToMessage() Message {
	return Message{
		Header: uint8(m.St),
		Params: m.Params,
	}
}

func NewResponseMessage(m Message) ResponseMessage {
	return ResponseMessage{
		St:     Status(m.Header),
		Params: m.Params,
	}
}
