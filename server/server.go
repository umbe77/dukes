// Copyright (c) 2023 Robeto Ughi
// 
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package server

import (
	"log"
	"net"

	"github.com/umbe77/dukes/cache"
	"github.com/umbe77/dukes/command"
	"github.com/umbe77/dukes/message"
)

type Server struct {
	ListenerAddr string
	ch           *cache.MemoryCache
	ln           net.Listener
}

func NewServer(addr string, cache *cache.MemoryCache) *Server {
	return &Server{
		ListenerAddr: addr,
		ch:           cache,
	}
}

func (s *Server) Close() {
	s.ln.Close()
}

func (s *Server) Run() error {
	var err error
	s.ln, err = net.Listen("tcp", s.ListenerAddr)
	if err != nil {
		return err
	}

	go func(ln net.Listener) {
		for {
			conn, err := s.ln.Accept()
			if err != nil {
				//TODO: logging
				log.Println(err)
				continue
			}
			go s.handleConnection(conn)
		}
	}(s.ln)
	return nil
}

func (s *Server) handleConnection(c net.Conn) {
	defer c.Close()
	for {
		m, err := message.Deserialize(c)
		if err != nil { //TODO: logging
			log.Printf("Connection closed by %s", c.RemoteAddr())
			break

		}
		//TODO: Logging message received (print ut the message)
		go s.handleCommand(message.NewRequestMessage(m), c)
	}
}

func (s *Server) handleCommand(m message.RequestMessage, c net.Conn) {
	var cmd command.Command
	switch m.Cmd {
	case message.CmdPing:
		cmd = &command.PingCommand{}
	case message.CmdSet:
		cmd = command.NewSetCommand(s.ch)
	case message.CmdGet:
		cmd = command.NewGetCommand(s.ch)
	case message.CmdHas:
		cmd = command.NewHasCommand(s.ch)
	case message.CmdDel:
		cmd = command.NewDelCommand(s.ch)
	case message.CmdDump:
		cmd = command.NewDumpCommand(s.ch)
	}

	for v := range cmd.Execute(m) {
		//TODO: Debug log response message
		// log.Println(v)
		c.Write(v)
	}
	log.Printf("Command %d, Handled\n", m.Cmd)
	//TODO: logging message executed
}
