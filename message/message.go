package message

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/umbe77/ucd/datatypes"
)

type Message struct {
	Header uint8
	Params []MessageParam
}

func (m Message) Serialize() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, m.Header)
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
		header  uint8
		readErr error
	)

	readErr = binary.Read(r, binary.LittleEndian, &header)
	if readErr != nil {
		return Message{}, readErr
	}
	var paramCount uint8
	readErr = binary.Read(r, binary.LittleEndian, &paramCount)
	if readErr != nil {
		return Message{}, readErr
	}

	m := Message{
		Header: header,
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
