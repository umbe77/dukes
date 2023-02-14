package command

import (
	"bytes"
	"encoding/binary"

	"github.com/umbe77/ucd/message"
)

type PingCommand struct {
}

func (c *PingCommand) Execute(m message.Message) <-chan []byte {
	ch := make(chan []byte)

	go func() {
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, message.BeginResp)
		binary.Write(buf, binary.LittleEndian, message.OK)
		binary.Write(buf, binary.LittleEndian, int32(len("Pong")))
		binary.Write(buf, binary.LittleEndian, []byte("Pong"))
		binary.Write(buf, binary.LittleEndian, message.EndResp)
		ch <- buf.Bytes()
		buf.Reset()
		close(ch)
	}()

	return ch

}
