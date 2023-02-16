package command

import (
	"github.com/umbe77/ucd/cache"
	"github.com/umbe77/ucd/datatypes"
	"github.com/umbe77/ucd/message"
)

type DumpCommand struct {
	mc *cache.MemoryCache
}

func NewDumpCommand(cache *cache.MemoryCache) *DumpCommand {
	return &DumpCommand{
		mc: cache,
	}
}

func (c *DumpCommand) Execute(m message.RequestMessage) <-chan []byte {
	ch := make(chan []byte)

	go func(m message.RequestMessage, mc *cache.MemoryCache) {

		for key := range mc.Dump() {
			ch <- message.ResponseMessage{
				St: message.OK,
				Params: []message.MessageParam{
					message.NewMessageParam(datatypes.String, key),
				},
			}.ToMessage().Serialize()
		}

		ch <- message.ResponseMessage{
			St:     message.EndResp,
			Params: []message.MessageParam{},
		}.ToMessage().Serialize()

		close(ch)
	}(m, c.mc)

	return ch
}
