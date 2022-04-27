package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	connection net.Conn
	name string
	room *room
	commands chan <- command
}

func (c *client) readInput() {
	for {
		message, err := bufio.NewReader(c.connection).ReadString('\n')

		if err != nil {
			return
		}

		if len(message) < 3 {
			continue
		}
			
		message = strings.Trim(message, "\r\n")
		arguments := strings.Split(message, " ")
		coms := strings.TrimSpace(arguments[0])

		switch coms {
		case "-nick" :
			c.commands <- command {
				id : comm_name,
				client: c,
				arguments: arguments,
			}
		case "-join" :
			c.commands <- command {
				id : comm_join,
				client: c,
				arguments: arguments,
			}
		case "-rooms" :
			c.commands <- command {
				id : comm_rooms,
				client: c,
				arguments: arguments,
			}
		case "-quit" :
			c.commands <- command {
				id : comm_quit,
				client: c,
				arguments: arguments,
			}
		case "-help" :
			c.commands <- command {
				id : comm_help,
				client : c,
				arguments : arguments,
			}
		case "-lvroom" :
			c.commands <- command {
				id : comm_lvroom,
				client: c,
				arguments: arguments,
			}
		case "-members" :
			c.commands <- command {
				id : comm_members,
				client: c,
				arguments: arguments,
			}
		default :{
				if arguments[0][0] != '-' {
					c.commands <- command {
						id : comm_message,
						client: c,
						arguments: arguments,
					}
				} else {
					c.err(fmt.Errorf("Unktown command : "+coms))
				}
			}
		}
	}
}

func (c *client) err(err error) {
	c.connection.Write([]byte("ERROR : " + err.Error() + "\n"))
}

func (c *client) message(message string) {
	c.connection.Write([]byte(">> " + message + "\n"))
}