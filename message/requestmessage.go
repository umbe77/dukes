package message

type RequestCommand uint8

const (
	CmdPing RequestCommand = iota
	CmdSet
	CmdGet
	CmdHas
	CmdDel
	CmdDump
)

type RequestMessage struct {
	Cmd    RequestCommand
	Params []MessageParam
}

func (m RequestMessage) ToMessage() Message {
	return Message{
		Header: uint8(m.Cmd),
		Params: m.Params,
	}
}

func NewRequestMessage(m Message) RequestMessage {
	return RequestMessage{
		Cmd:    RequestCommand(m.Header),
		Params: m.Params,
	}
}
