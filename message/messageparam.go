package message

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/umbe77/dukes/datatypes"
)

type MessageParam struct {
	Kind  datatypes.DataType
	Value []byte
	Len   int32
}

func (m MessageParam) ToAny() any {
	switch m.Kind {
	case datatypes.Int:
		var v int32
		b := bytes.NewReader(m.Value)
		binary.Read(b, binary.LittleEndian, &v)
		return v
	case datatypes.String:
		return string(m.Value)
	case datatypes.Bool:
		return m.Value[0] == 1
	case datatypes.Date:
		var unix int64
		b := bytes.NewReader(m.Value)
		binary.Read(b, binary.LittleEndian, &unix)
		return time.Unix(unix, 0)
	}
	return nil
}

func NewMessageParam(kind datatypes.DataType, value any) MessageParam {
	p := MessageParam{
		Kind: kind,
	}
	switch kind {
	case datatypes.Int:
		buff := new(bytes.Buffer)
		binary.Write(buff, binary.LittleEndian, value.(int32))
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
