package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("new client is connected %s", conn.RemoteAddr().String())
	c := &client{
		conn:     conn,
		name:     "anonymous",
		commands: s.commands,
	}
	c.readInput()
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NAME:
			s.name(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client, cmd.args)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
		}
	}
}

func (s *server) name(c *client, args []string) {
	if len(args) < 2 {
		c.clientMsg("nick is required")
		return
	}

	c.name = args[1]
	c.clientMsg(fmt.Sprintf("all right, I will call you %s", c.name))
}
func (s *server) join(c *client, args []string) {
	if len(args) < 2 {
		c.clientMsg("room name is required")
		return
	}

	roomName := args[1]

	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			roomName: roomName,
			members:  make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.name))

	c.clientMsg(fmt.Sprintf("welcome to %s", roomName))
}
func (s *server) listRooms(c *client, args []string) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.clientMsg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")))
}
func (s *server) msg(c *client, args []string) {
	if len(args) < 2 {
		c.clientMsg("message is required")
		return
	}

	msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, c.name+": "+msg)
}
func (s *server) quit(c *client, args []string) {
	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.clientMsg("sad to see you go !!")
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.roomName]
		delete(s.rooms[c.room.roomName].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.name))
	}
}
