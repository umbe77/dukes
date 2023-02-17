// Copyright (c) 2023 Robeto Ughi
// 
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package command

import (
	"github.com/umbe77/dukes/cache"
	"github.com/umbe77/dukes/datatypes"
	"github.com/umbe77/dukes/message"
)

type SetCommand struct {
	mc *cache.MemoryCache
}

func NewSetCommand(c *cache.MemoryCache) *SetCommand {
	return &SetCommand{
		mc: c,
	}
}

func getSetResp(m message.RequestMessage, c *cache.MemoryCache) message.ResponseMessage {
	if len(m.Params) != 2 {
		return message.ResponseMessage{
			St: message.BadFormat,
			Params: []message.MessageParam{
				message.NewMessageParam(datatypes.String, "Set message should have 2 params, key and value"),
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
	value := &cache.CacheValue{
		Kind:  m.Params[1].Kind,
		Value: m.Params[1].ToAny(),
	}
	if err := c.Set(key, value); err != nil {
		return message.ResponseMessage{
			St: message.Error,
			Params: []message.MessageParam{
				message.NewMessageParam(datatypes.String, err.Error()),
			},
		}
	}
	return message.ResponseMessage{
		St: message.OK,
		Params: []message.MessageParam{
			message.NewMessageParam(datatypes.String, key),
		},
	}
}

// TODO: MAKE TEST FOR Set Command
func (c *SetCommand) Execute(m message.RequestMessage) <-chan []byte {
	ch := make(chan []byte)

	go func(m message.RequestMessage, mc *cache.MemoryCache) {
		ch <- getSetResp(m, mc).ToMessage().Serialize()

		ch <- message.ResponseMessage{
			St:     message.EndResp,
			Params: []message.MessageParam{},
		}.ToMessage().Serialize()

		close(ch)
	}(m, c.mc)

	return ch
}
