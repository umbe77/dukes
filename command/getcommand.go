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

type GetCommand struct {
	mc *cache.Cache
}

func NewGetCommand(cache *cache.Cache) *GetCommand {
	return &GetCommand{
		mc: cache,
	}
}

func getGetResp(m message.RequestMessage, c *cache.Cache) message.ResponseMessage {
	if len(m.Params) != 1 {
		return message.ResponseMessage{
			St: message.BadFormat,
			Params: []message.MessageParam{
				message.NewMessageParam(datatypes.String, "Get message should have 1 params, key"),
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
	cacheValue, err := c.Get(key)
	if err != nil {
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
			message.NewMessageParam(cacheValue.Kind, cacheValue.Value),
		},
	}
}

// TODO: MAKE TEST
func (c *GetCommand) Execute(m message.RequestMessage) <-chan []byte {
	ch := make(chan []byte)

	go func(m message.RequestMessage, mc *cache.Cache) {
		ch <- getGetResp(m, mc).ToMessage().Serialize()

		ch <- message.ResponseMessage{
			St:     message.EndResp,
			Params: []message.MessageParam{},
		}.ToMessage().Serialize()

		close(ch)
	}(m, c.mc)

	return ch
}
