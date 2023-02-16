package command

import "github.com/umbe77/dukes/message"

type Command interface {
	Execute(message.RequestMessage) <-chan []byte
}
