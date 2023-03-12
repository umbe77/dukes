package command

import (
	"time"

	"github.com/hashicorp/raft"

	"github.com/umbe77/dukes/datatypes"
	"github.com/umbe77/dukes/message"
)

type DelCommand struct {
	ra *raft.Raft
}

func NewDelCommand(ra *raft.Raft) *DelCommand {
	return &DelCommand{
		ra: ra,
	}
}

func getDelResp(m message.RequestMessage, ra *raft.Raft) message.ResponseMessage {
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
	err := ra.Apply(m.ToMessage().Serialize(), time.Millisecond*50).Error()
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

	go func(m message.RequestMessage, ra *raft.Raft) {

		ch <- getDelResp(m, ra).ToMessage().Serialize()

		ch <- message.ResponseMessage{
			St:     message.EndResp,
			Params: []message.MessageParam{},
		}.ToMessage().Serialize()

		close(ch)
	}(m, c.ra)

	return ch
}
