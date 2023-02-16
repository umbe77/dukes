package command

import (
	"github.com/umbe77/dukes/datatypes"
	"github.com/umbe77/dukes/message"
)

type PingCommand struct {
}

func (c *PingCommand) Execute(m message.RequestMessage) <-chan []byte {
	ch := make(chan []byte)

	go func() {

		ch <- message.ResponseMessage{
			St: message.OK,
			Params: []message.MessageParam{
				message.NewMessageParam(datatypes.String, "Pong"),
			},
		}.ToMessage().Serialize()

		ch <- message.ResponseMessage{
			St:     message.EndResp,
			Params: []message.MessageParam{},
		}.ToMessage().Serialize()

		close(ch)
	}()

	return ch

}
