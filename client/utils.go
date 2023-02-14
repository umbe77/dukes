package client

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/umbe77/ucd/message"
)

func processSimpleResopnse(conn net.Conn) (string, error) {
	//First Read from response should be a Status BeginResp
	var bStatus message.Status
	err := binary.Read(conn, binary.LittleEndian, &bStatus)
	if err != nil {
		return "", err
	}
	if bStatus != message.BeginResp {
		return "", fmt.Errorf("response message bad formatted: %v, should be: %v", bStatus, message.BeginResp)
	}

	//Second and third packet shuold be a status ok and message Pong
	status := make([]byte, 1)
	err = binary.Read(conn, binary.LittleEndian, &status)
	if err != nil {
		return "", err
	}
	var respLen int32
	err = binary.Read(conn, binary.LittleEndian, &respLen)
	if err != nil {
		return "", err
	}

	resp := make([]byte, respLen)
	err = binary.Read(conn, binary.LittleEndian, &resp)
	if err != nil {
		return "", err
	}

	//Last Packet should be a EndResp status
	var eStatus message.Status
	err = binary.Read(conn, binary.LittleEndian, &eStatus)
	if err != nil {
		return "", err
	}
	if eStatus != message.EndResp {
		return "", fmt.Errorf("response message bad formatted: %v, should be: %v", eStatus, message.EndResp)
	}

	if message.Status(status[0]) != message.OK {
		return "", fmt.Errorf("error: %s", resp)
	}
	return string(resp), nil

}
