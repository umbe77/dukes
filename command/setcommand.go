package command

import (
	"bytes"
	"encoding/binary"

	"github.com/umbe77/ucd/cache"
	"github.com/umbe77/ucd/datatypes"
	"github.com/umbe77/ucd/message"
)

type SetCommand struct {
	c *cache.MemoryCache
}

func NewSetCommand(c *cache.MemoryCache) *SetCommand {
	return &SetCommand{
		c: c,
	}
}

// TODO: MAKE TEST FOR Set Command
func (c *SetCommand) Execute(m message.Message) <-chan []byte {
	ch := make(chan []byte)

	go func() {
		status := message.BadFormat
		msg := []byte("Set message should have 2 params, key and value")
		if len(m.Params) == 2 {
			if m.Params[0].Kind == datatypes.String {
				key := string(m.Params[0].Value)
				value := &cache.CacheValue{
					Kind:  m.Params[1].Kind,
					Value: m.Params[1].Value,
				}
				status = message.OK
				msg = m.Params[0].Value
				if err := c.c.Set(key, value); err != nil {
					status = message.Error
					msg = []byte(err.Error())
				}
			}
		}
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, message.BeginResp)
		binary.Write(buf, binary.LittleEndian, status)
		binary.Write(buf, binary.LittleEndian, int32(len(msg)))
		binary.Write(buf, binary.LittleEndian, msg)
		binary.Write(buf, binary.LittleEndian, message.EndResp)
		ch <- buf.Bytes()
		close(ch)
	}()

	return ch
}
