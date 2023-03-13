package command

import (
	"fmt"

	"github.com/hashicorp/raft"

	"github.com/umbe77/dukes/datatypes"
	"github.com/umbe77/dukes/message"
)

type JoinCommand struct {
	ra *raft.Raft
}

func NewJoinCommand(ra *raft.Raft) *JoinCommand {
	return &JoinCommand{
		ra: ra,
	}
}

func join(id, addr string, ra *raft.Raft) error {
	config := ra.GetConfiguration()
	if err := config.Error(); err != nil {
		//TODO: Add logging
		return err
	}

	nodeId := raft.ServerID(id)
	nodeAddr := raft.ServerAddress(addr)

	for _, srv := range config.Configuration().Servers {
		srvAddr := srv.Address
		srvNodeId := srv.ID
		if srvAddr == nodeAddr && srvNodeId == nodeId {
			return nil
		}
		if srvAddr == nodeAddr || srvNodeId == nodeId {
			if err := ra.RemoveServer(nodeId, 0, 0).Error(); err != nil {
				return fmt.Errorf("error removing node %s at %s: %s", id, addr, err)
			}
		}
	}

	f := ra.AddVoter(nodeId, nodeAddr, 0, 0)
	if err := f.Error(); err != nil {
		return err
	}

	return nil
}

func getJoinResp(m message.RequestMessage, ra *raft.Raft) message.ResponseMessage {
	if len(m.Params) != 2 {
		return message.ResponseMessage{
			St: message.BadFormat,
			Params: []message.MessageParam{
				message.NewMessageParam(datatypes.String, "Join message should have 2 params, key and value"),
			},
		}
	}

	if m.Params[0].Kind != datatypes.String {
		return message.ResponseMessage{
			St: message.BadFormat,
			Params: []message.MessageParam{
				message.NewMessageParam(datatypes.String, "NodeId params must be a string kind"),
			},
		}
	}
	if m.Params[1].Kind != datatypes.String {
		return message.ResponseMessage{
			St: message.BadFormat,
			Params: []message.MessageParam{
				message.NewMessageParam(datatypes.String, "NodeAddres params must be a string kind"),
			},
		}
	}

	nodeId := string(m.Params[0].Value)
	nodeAddr := string(m.Params[1].Value)

	if err := join(nodeId, nodeAddr, ra); err != nil {
		return message.ResponseMessage{
			St: message.Error,
			Params: []message.MessageParam{
				message.NewMessageParam(datatypes.String, fmt.Sprintf("cannot join %s at %s: %s", nodeId, nodeAddr, err)),
			},
		}
	}
	return message.ResponseMessage{
		St:     message.OK,
		Params: []message.MessageParam{},
	}
}

func (c *JoinCommand) Execute(m message.RequestMessage) <-chan []byte {
	ch := make(chan []byte)

	go func(m message.RequestMessage, ra *raft.Raft) {

		ch <- getJoinResp(m, ra).ToMessage().Serialize()

		ch <- message.ResponseMessage{
			St:     message.EndResp,
			Params: []message.MessageParam{},
		}.ToMessage().Serialize()

		close(ch)

	}(m, c.ra)

	return ch
}
