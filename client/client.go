package client

import (
	"net"

	"github.com/umbe77/dukes/datatypes"
	"github.com/umbe77/dukes/message"
)

type Client struct {
	ServerAddress string
	conn          net.Conn
}

func NewClient(serverAddr string) (*Client, error) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil, err
	}

	return &Client{
		ServerAddress: serverAddr,
		conn:          conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Ping() (string, error) {
	pingMsg := message.RequestMessage{
		Cmd:    message.CmdPing,
		Params: make([]message.MessageParam, 0),
	}

	req := pingMsg.ToMessage()
	_, sendErr := c.conn.Write(req.Serialize())
	if sendErr != nil {
		return "", sendErr
	}

	return precessSimpleResponse(c.conn)
}

func (c *Client) Join(nodeId, nodeAddr string) error {
	joinMsg := message.RequestMessage{
		Cmd: message.CmdJoin,
		Params: []message.MessageParam{
			message.NewMessageParam(datatypes.String, nodeId),
			message.NewMessageParam(datatypes.String, nodeAddr),
		},
	}.ToMessage()

	if _, err := c.conn.Write(joinMsg.Serialize()); err != nil {
		return err
	}

	return processVoidResponse(c.conn)

}

func (c *Client) Set(key string, kind datatypes.DataType, value any) (string, error) {
	setMsg := message.RequestMessage{
		Cmd: message.CmdSet,
		Params: []message.MessageParam{
			message.NewMessageParam(datatypes.String, key),
			message.NewMessageParam(kind, value),
		},
	}.ToMessage()

	if _, err := c.conn.Write(setMsg.Serialize()); err != nil {
		return "", err
	}

	return precessSimpleResponse(c.conn)
}

func (c *Client) Get(key string) (message.MessageParam, error) {
	getMsg := message.RequestMessage{
		Cmd: message.CmdGet,
		Params: []message.MessageParam{
			message.NewMessageParam(datatypes.String, key),
		},
	}

	if _, err := c.conn.Write(getMsg.ToMessage().Serialize()); err != nil {
		return message.MessageParam{}, err
	}

	return processGetResponse(c.conn)
}

func (c *Client) Has(key string) (bool, error) {
	hasMsg := message.RequestMessage{
		Cmd: message.CmdHas,
		Params: []message.MessageParam{
			message.NewMessageParam(datatypes.String, key),
		},
	}
	if _, err := c.conn.Write(hasMsg.ToMessage().Serialize()); err != nil {
		return false, err
	}

	return processHasResponse(c.conn)
}

func (c *Client) Del(key string) error {
	delMsg := message.RequestMessage{
		Cmd: message.CmdDel,
		Params: []message.MessageParam{
			message.NewMessageParam(datatypes.String, key),
		},
	}
	if _, err := c.conn.Write(delMsg.ToMessage().Serialize()); err != nil {
		return err
	}

	return processVoidResponse(c.conn)
}

func (c *Client) Dump() <-chan string {
	dumpMsg := message.RequestMessage{
		Cmd:    message.CmdDump,
		Params: []message.MessageParam{},
	}

	if _, err := c.conn.Write(dumpMsg.ToMessage().Serialize()); err != nil {
		return nil
	}

	return processDumpResponse(c.conn)

}
