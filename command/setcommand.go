package command

import (
	"time"

	"github.com/hashicorp/raft"

	"github.com/umbe77/dukes/datatypes"
	"github.com/umbe77/dukes/message"
)

type SetCommand struct {
	ra *raft.Raft
}

func NewSetCommand(ra *raft.Raft) *SetCommand {
	return &SetCommand{
		ra: ra,
	}
}

func getSetResp(m message.RequestMessage, ra *raft.Raft) message.ResponseMessage {
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

	if err := ra.Apply(m.ToMessage().Serialize(), time.Millisecond*50).Error(); err != nil {
		return message.ResponseMessage{
			St: message.Error,
			Params: []message.MessageParam{
				message.NewMessageParam(datatypes.String, err.Error()),
			},
		}
	}

	key := string(m.Params[0].Value)
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

	go func(m message.RequestMessage, ra *raft.Raft) {
		ch <- getSetResp(m, ra).ToMessage().Serialize()

		ch <- message.ResponseMessage{
			St:     message.EndResp,
			Params: []message.MessageParam{},
		}.ToMessage().Serialize()

		close(ch)
	}(m, c.ra)

	return ch
}
