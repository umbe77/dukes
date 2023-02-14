package message

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"

	"github.com/umbe77/ucd/datatypes"
)

type Status uint8

const (
	BeginResp Status = iota
	EndResp
	OK
	BadFormat
	Error
	Nope
)

func (s Status) Bytes() []byte {
	return []byte{uint8(s)}
}

type Command uint8

const (
	CmdPing Command = iota
	CmdSet
	CmdGet
	CmdDel
	CmdDump
)

type MessageParam struct {
	Kind  datatypes.DataType
	Value []byte
	Len   int32
}

func NewMessageParam(kind datatypes.DataType, value any) MessageParam {
	p := MessageParam{
		Kind: kind,
	}
	switch kind {
	case datatypes.Int:
		buff := new(bytes.Buffer)
		binary.Write(buff, binary.LittleEndian, int32(value.(int)))
		buffBytes := buff.Bytes()
		p.Value = buffBytes
		p.Len = int32(len(buffBytes))
	case datatypes.String:
		p.Value = []byte(value.(string))
		p.Len = int32(len(value.(string)))
	case datatypes.Bool:
		p.Value = []byte{0}
		if value.(bool) {
			p.Value = []byte{1}
		}
		p.Len = 1
	case datatypes.Date:
		buff := new(bytes.Buffer)
		binary.Write(buff, binary.LittleEndian, value.(time.Time).Unix())
		buffBytes := buff.Bytes()
		p.Value = buffBytes
		p.Len = int32(len(buffBytes))
	}

	return p
}

type Message struct {
	Cmd    Command
	Params []MessageParam
}

func (m *Message) Serialize() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, m.Cmd)
	paramCount := uint8(len(m.Params))
	binary.Write(buf, binary.LittleEndian, paramCount)
	for _, p := range m.Params {
		binary.Write(buf, binary.LittleEndian, p.Kind)
		binary.Write(buf, binary.LittleEndian, p.Len)
		binary.Write(buf, binary.LittleEndian, p.Value)
	}

	return buf.Bytes()
}

func Deserialize(r io.Reader) (Message, error) {
	var (
		cmd     Command
		readErr error
	)

	readErr = binary.Read(r, binary.LittleEndian, &cmd)
	if readErr != nil {
		return Message{}, readErr
	}
	var paramCount uint8
	readErr = binary.Read(r, binary.LittleEndian, &paramCount)
	if readErr != nil {
		return Message{}, readErr
	}

	m := Message{
		Cmd:    cmd,
		Params: make([]MessageParam, paramCount),
	}
	for i := 0; i < int(paramCount); i++ {
		var pKind datatypes.DataType
		readErr = binary.Read(r, binary.LittleEndian, &pKind)
		if readErr != nil {
			return Message{}, readErr
		}
		var pLen int32
		readErr = binary.Read(r, binary.LittleEndian, &pLen)
		if readErr != nil {
			return Message{}, readErr
		}
		pVal := make([]byte, pLen)
		readErr = binary.Read(r, binary.LittleEndian, &pVal)
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return Message{}, readErr
		}
		m.Params[i] = MessageParam{
			Kind:  pKind,
			Len:   pLen,
			Value: pVal,
		}
	}
	return m, nil
}
