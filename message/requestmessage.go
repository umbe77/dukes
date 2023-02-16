package message

type Command uint8

const (
	CmdPing Command = iota
	CmdSet
	CmdGet
	CmdDel
	CmdDump
)

type RequestMessage struct {
	Cmd    Command
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
		Cmd:    Command(m.Header),
		Params: m.Params,
	}
}
