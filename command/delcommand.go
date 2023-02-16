package command

import (
	"github.com/umbe77/ucd/cache"
	"github.com/umbe77/ucd/datatypes"
	"github.com/umbe77/ucd/message"
)

type DelCommand struct {
	mc *cache.MemoryCache
}

func NewDelCommand(cache *cache.MemoryCache) *DelCommand {
	return &DelCommand{
		mc: cache,
	}
}

func getDelResp(m message.RequestMessage, c *cache.MemoryCache) message.ResponseMessage {
	if len(m.Params) != 1 {
		return message.ResponseMessage{
			St: message.BadFormat,
			Params: []message.MessageParam{
				message.NewMessageParam(datatypes.String, "Del message should have 1 params, key"),
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
	err := c.Del(key)
	if err != nil {
		return message.ResponseMessage{
			St: message.Error,
			Params: []message.MessageParam{
				message.NewMessageParam(datatypes.String, err.Error()),
			},
		}
	}
	return message.ResponseMessage{
		St:     message.OK,
		Params: []message.MessageParam{},
	}
}
func (c *DelCommand) Execute(m message.RequestMessage) <-chan []byte {
	ch := make(chan []byte)

	go func(m message.RequestMessage, mc *cache.MemoryCache) {

		ch <- getDelResp(m, mc).ToMessage().Serialize()

		ch <- message.ResponseMessage{
			St:     message.EndResp,
			Params: []message.MessageParam{},
		}.ToMessage().Serialize()

		close(ch)
	}(m, c.mc)

	return ch
}
