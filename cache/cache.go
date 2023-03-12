// Copyright (c) 2023 Robeto Ughi
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cache

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/hashicorp/raft"

	"github.com/umbe77/dukes/datatypes"
	"github.com/umbe77/dukes/message"
)

type CacheValue struct {
	Kind  datatypes.DataType
	Value any
}

type Cache struct {
	sync.RWMutex
	c map[string]*CacheValue
}

func NewCache() *Cache {
	return &Cache{
		c: make(map[string]*CacheValue),
	}
}

func (c *Cache) Set(key string, value *CacheValue) error {
	c.Lock()
	c.c[key] = value
	c.Unlock()
	return nil
}

func (c *Cache) Has(key string) bool {
	var result bool
	c.RLock()
	_, result = c.c[key]
	c.RUnlock()
	return result
}

func (c *Cache) Get(key string) (*CacheValue, error) {
	var (
		v  *CacheValue
		ok bool
	)

	c.RLock()
	if v, ok = c.c[key]; !ok {
		return nil, fmt.Errorf("key %s not present in cache", key)
	}
	c.RUnlock()
	return v, nil
}

func (c *Cache) Del(key string) error {

	c.Lock()
	delete(c.c, key)
	c.Unlock()
	return nil
}

func (c *Cache) Dump() <-chan string {
	keysCh := make(chan string)

	go func(mc *Cache) {

		for k := range c.c {
			keysCh <- k
		}

		close(keysCh)
	}(c)

	return keysCh
}

func applyError(err error) message.ResponseMessage {
	return message.ResponseMessage{
		St: message.Error,
		Params: []message.MessageParam{
			message.NewMessageParam(datatypes.String, err.Error()),
		},
	}
}

func (c *Cache) Apply(log *raft.Log) any {
	switch log.Type {
	case raft.LogCommand:
		reader := bytes.NewReader(log.Data)
		msg, err := message.Deserialize(reader)
		if err != nil {
			return applyError(err)
		}
		req := message.NewRequestMessage(msg)

		switch req.Cmd {
		case message.CmdSet:
			value := &CacheValue{
				Kind:  msg.Params[1].Kind,
				Value: msg.Params[1].ToAny(),
			}
			err := c.Set(string(req.Params[0].Value), value)
			if err != nil {
				return applyError(err)
			}
			return message.ResponseMessage{
				St: message.OK,
				Params: []message.MessageParam{
					msg.Params[1],
				},
			}
		case message.CmdDel:
			key := string(req.Params[0].Value)
			err := c.Del(key)
			if err != nil {
				return applyError(err)
			}
			return message.ResponseMessage{
				St:     message.OK,
				Params: []message.MessageParam{},
			}
		}
	}
	return nil
}

func (c *Cache) Snapshot() (raft.FSMSnapshot, error) {
	return &fsmSnapshot{}, nil
}

func (c *Cache) Restore(snapshot io.ReadCloser) error {
	return nil
}

type fsmSnapshot struct{}

func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error { return nil }

func (f *fsmSnapshot) Release() {}
