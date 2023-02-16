package server

import (
	"log"
	"net"

	"github.com/umbe77/ucd/cache"
	"github.com/umbe77/ucd/command"
	"github.com/umbe77/ucd/message"
)

type Server struct {
	listenerAddr string
	ch           *cache.MemoryCache
	ln           net.Listener
}

func NewServer(addr string, cache *cache.MemoryCache) *Server {
	return &Server{
		listenerAddr: addr,
		ch:           cache,
	}
}

func (s *Server) Run() error {
	var err error
	s.ln, err = net.Listen("tcp", s.listenerAddr)
	defer s.ln.Close()
	if err != nil {
		return err
	}

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			//TODO: logging
			log.Println(err)
			continue
		}
		go s.handleConnection(conn)
	}
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
	}

	for v := range cmd.Execute(m) {
		//TODO: Debug log response message
		// log.Println(v)
		c.Write(v)
	}
	log.Printf("Command %d, Handled\n", m.Cmd)
	//TODO: logging message executed
}
