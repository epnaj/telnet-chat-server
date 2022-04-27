package main

import (
	"net"
)

type room struct {
	name string
	members map[net.Addr] *client
}

func (r *room) broadcast(sender *client, message string) {
	for address, mem := range r.members {
		if address != sender.connection.RemoteAddr() {
			mem.message(message)
		}
	}
}