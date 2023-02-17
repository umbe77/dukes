// Copyright (c) 2023 Robeto Ughi
// 
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package command

import "github.com/umbe77/dukes/message"

type Command interface {
	Execute(message.RequestMessage) <-chan []byte
}
