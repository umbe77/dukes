package client

import (
	"fmt"
	"net"

	"github.com/umbe77/ucd/datatypes"
	"github.com/umbe77/ucd/message"
)

func GetResponseMessage(conn net.Conn) (message.ResponseMessage, error) {
	m, err := message.Deserialize(conn)
	if err != nil {
		return message.ResponseMessage{}, err
	}
	return message.NewResponseMessage(m), nil
}

func precessSimpleResponse(conn net.Conn) (string, error) {

	var resp string
	var respError error
	msgReceivedCount := 0
	for {
		msgReceivedCount++
		m, err := GetResponseMessage(conn)
		if err != nil {
			return "", err
		}
		if m.St == message.EndResp {
			break
		}
		if msgReceivedCount == 1 {
			if len(m.Params) != 1 {
				respError = fmt.Errorf("bad format, not enough params")
				continue
			}
			resp = string(m.Params[0].Value)
			if m.St == message.Error {
				respError = fmt.Errorf(string(m.Params[0].Value))
			}
		}
	}

	return resp, respError

}

func processGetResponse(conn net.Conn) (message.MessageParam, error) {

	var resp message.MessageParam
	var respError error
	msgReceivedCount := 0
	for {
		msgReceivedCount++
		m, err := GetResponseMessage(conn)
		if err != nil {
			return message.MessageParam{}, err
		}
		if m.St == message.EndResp {
			break
		}
		if msgReceivedCount == 1 {
			if len(m.Params) != 1 {
				respError = fmt.Errorf("bad format, not enough params")
				continue
			}
			resp = m.Params[0]
			if m.St == message.Error {
				respError = fmt.Errorf(string(m.Params[0].Value))
			}
		}
	}

	return resp, respError

}

func processHasResponse(conn net.Conn) (bool, error) {
	var resp bool
	var respError error
	msgReceivedCount := 0
	for {
		msgReceivedCount++
		m, err := GetResponseMessage(conn)
		if err != nil {
			return false, err
		}
		if m.St == message.EndResp {
			break
		}
		if msgReceivedCount == 1 {
			if len(m.Params) != 1 {
				respError = fmt.Errorf("bad format, not enough params")
				continue
			}
			respParam := m.Params[0]
			if m.St == message.Error {
				respError = fmt.Errorf(string(m.Params[0].Value))
			}
			if respParam.Kind != datatypes.Bool {
				respError = fmt.Errorf("bad format, message not in bool format")
			}
			resp = respParam.ToAny().(bool)
		}
	}

	return resp, respError
}

func processDelResponse(conn net.Conn) error {
	var respError error
	msgReceivedCount := 0
	for {
		msgReceivedCount++
		m, err := GetResponseMessage(conn)
		if err != nil {
			return err
		}
		if m.St == message.EndResp {
			break
		}
		if msgReceivedCount == 1 {
			if m.St == message.Error {
				respError = fmt.Errorf(string(m.Params[0].Value))
			}
		}
	}

	return respError
}
