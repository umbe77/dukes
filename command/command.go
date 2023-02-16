package command

import "github.com/umbe77/ucd/message"

type Command interface {
	Execute(message.RequestMessage) <-chan []byte
}
