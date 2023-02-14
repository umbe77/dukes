package client

import (
	"net"

	"github.com/umbe77/ucd/datatypes"
	"github.com/umbe77/ucd/message"
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
	pingMsg := message.Message{
		Cmd:    message.CmdPing,
		Params: make([]message.MessageParam, 0),
	}

	_, sendErr := c.conn.Write(pingMsg.Serialize())
	if sendErr != nil {
		return "", sendErr
	}

	return processSimpleResopnse(c.conn)

}

func (c *Client) Set(key string, kind datatypes.DataType, value any) (string, error) {
	setMsg := message.Message{
		Cmd: message.CmdSet,
		Params: []message.MessageParam{
			message.NewMessageParam(datatypes.String, key),
			message.NewMessageParam(kind, value),
		},
	}

	if _, err := c.conn.Write(setMsg.Serialize()); err != nil {
		return "", err
	}

	return processSimpleResopnse(c.conn)
}

func (c *Client) Get(key string) (message.MessageParam, error) {
	getMsg := message.Message{
		Cmd: message.CmdGet,
		Params: []message.MessageParam{
			message.NewMessageParam(datatypes.String, key),
		},
	}

	if _, err := c.conn.Write(getMsg.Serialize()); err != nil {
		return message.MessageParam{}, err
	}
	return message.MessageParam{}, nil
}
