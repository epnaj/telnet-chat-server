package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func NewServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for command := range s.commands {
		switch command.id {
		case comm_name : 
			s.name(command.client, command.arguments)
		case comm_join : 
			s.join(command.client, command.arguments)
		case comm_message : 
			s.message(command.client, command.arguments)
		case comm_quit : 
			s.quit(command.client, command.arguments)
		case comm_rooms : 
			s.roomsFunctionExept(command.client, command.arguments)
		case comm_help :
			s.helpDisplay(command.client)
		case comm_lvroom :
			s.leaveRoom(command.client)
		case comm_members :
			s.membersDisplay(command.client)
		}
	}
}

func (s *server) newClient(connection net.Conn) {
	log.Println("New client has connected :", connection.RemoteAddr().String())

	c := &client {
		connection : connection,
		name : "anonymus",
		commands : s.commands,
	}
	s.helpDisplay(c)
	c.readInput()
}

func (s *server) name(c *client, arguments []string){
	
	if len(arguments) < 2 {
		c.message("Name can't be empty")
		return
	}
	if len(arguments[1]) > 0 && len(strings.ReplaceAll(arguments[1], " ", "")) > 0 {
		joinedArguments := strings.Join(arguments[1:], " ")
		c.message(fmt.Sprint("Name changed to '", joinedArguments, "'"))
		if c.room != nil {
			c.room.broadcast(c, "'" + c.name + "' changed nick to '" + joinedArguments + "'")
		}
		c.name = joinedArguments
	} else {
		c.message("Name can't be empty")
	}
}

func (s *server) join(c *client, arguments []string){
	if len(arguments) < 2 {
		c.message("You must enter ther room name")
		return
	}
	if len(strings.ReplaceAll(arguments[1], " ", "")) == 0 {
		c.message("You must enter ther room name")
		return
	}

	roomName :=  strings.Join(arguments[1:], " ")
	r, exists := s.rooms[roomName]

	if !exists {
		r = &room{
			name : roomName,
			members : make(map[net.Addr] *client),
		}
		s.rooms[roomName] = r
	} else if c.room == r {
		c.message("You already are in '" + roomName + "'")
		return
	}

	r.members[c.connection.RemoteAddr()] = c

	s.quitCurrentRoom(c)

	c.room = r

	r.broadcast(c, fmt.Sprint(c.name, " has joined the room"))
	
	c.message(fmt.Sprint("welcome to ", r.name))
}

func (s *server) message(c *client, arguments []string){
	if c.room == nil {
		// errors shouldn't be capitalised xd
		c.err(errors.New("you must join the room first"))
		return
	}
	c.room.broadcast(c, c.name + " : " + strings.Join(arguments[0:], " "))
}

func (s *server) quit(c *client, arguments []string){
	log.Println("Client has disconnected : ", c.connection.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.message("Bye")
	c.connection.Close()
}

func (s *server) roomsFunctionExept(c *client, arguments []string){
	var rooms [] string
	for name := range s.rooms {
		rooms = append(rooms, name)	 
	}
	if len(rooms) > 0 {
		c.message(fmt.Sprint("avaliable rooms are : ", strings.Join(rooms, ", ")))
	} else {
		c.message("there are no open rooms, go ahead and acreate one!")
	}
}

func (s *server) helpDisplay(c *client) {
	c.message(HELP_DISPLAY_STRING)
}

func (s *server) quitCurrentRoom(c *client) {	
	if c.room != nil {
		delete(c.room.members, c.connection.RemoteAddr())
		c.room.broadcast(c, fmt.Sprint(c.name, " has left the room"))
	}

	if len(s.rooms) > 0 && c.room != nil {
		if len(s.rooms[c.room.name].members) == 0 {
			delete(s.rooms, c.room.name)
		}
	}
}

func (s *server) leaveRoom(c *client) {
	if c.room == nil {
		c.message("If you want to leave server use '-quit'")
		return
	}
	s.quitCurrentRoom(c)
	
	c.room = nil
}

func (s *server) membersDisplay(c *client) {
	if c.room != nil {
		rmPoint := s.rooms[c.room.name]
		c.message("Members :")
		for addr := range rmPoint.members {
			c.message(rmPoint.members[addr].name)
		}
	} else {
		c.message("You must be in the room")
	}
}