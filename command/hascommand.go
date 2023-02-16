package command

import (
	"github.com/umbe77/dukes/cache"
	"github.com/umbe77/dukes/datatypes"
	"github.com/umbe77/dukes/message"
)

type HasCommand struct {
	mc *cache.MemoryCache
}

func NewHasCommand(cache *cache.MemoryCache) *HasCommand {
	return &HasCommand{
		mc: cache,
	}
}

func getHasResp(m message.RequestMessage, c *cache.MemoryCache) message.ResponseMessage {
	if len(m.Params) != 1 {
		return message.ResponseMessage{
			St: message.BadFormat,
			Params: []message.MessageParam{
				message.NewMessageParam(datatypes.String, "Has message should have 1 params, key"),
			},
		}
	}
	if m.Params[0].Kind != datatypes.String {
		return message.ResponseMessage{
			St: message.BadFormat,
			Params: []message.MessageParam{
				message.NewMessageParam(datatypes.String, "Key params must be a string kind"),
			},
		}
	}
	key := string(m.Params[0].Value)
	hasKey := c.Has(key)

	return message.ResponseMessage{
		St: message.OK,
		Params: []message.MessageParam{
			message.NewMessageParam(datatypes.Bool, hasKey),
		},
	}
}

// TODO: MAKE TEST
func (c *HasCommand) Execute(m message.RequestMessage) <-chan []byte {
	ch := make(chan []byte)

	go func(m message.RequestMessage, mc *cache.MemoryCache) {

		ch <- getHasResp(m, mc).ToMessage().Serialize()

		ch <- message.ResponseMessage{
			St:     message.EndResp,
			Params: []message.MessageParam{},
		}.ToMessage().Serialize()

		close(ch)
	}(m, c.mc)

	return ch
}
