package command

import (
	"github.com/umbe77/dukes/cache"
	"github.com/umbe77/dukes/datatypes"
	"github.com/umbe77/dukes/message"
)

type DumpCommand struct {
	mc *cache.Cache
}

func NewDumpCommand(cache *cache.Cache) *DumpCommand {
	return &DumpCommand{
		mc: cache,
	}
}

func (c *DumpCommand) Execute(m message.RequestMessage) <-chan []byte {
	ch := make(chan []byte)

	go func(m message.RequestMessage, mc *cache.Cache) {

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
