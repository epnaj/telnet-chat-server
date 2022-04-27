package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	port := ":8888"
	if len(os.Args) > 1{
		n, err := strconv.Atoi(os.Args[1]) 
		if err != nil {
			fmt.Println("Using default port", port, "due to error", err)
		} else {
			port = ":" + strings.TrimLeft(os.Args[1], "0")
		}
		n++
	} else {
		fmt.Println("Using default port", port)
	}
	
	s := NewServer()
	go s.run()
	listener, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalln("Server has not started", err.Error())
	}
	
	defer listener.Close()

	log.Println("Server is up at", port)

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Println("Connection can't be accepted for some reason", err.Error())
		} else {
			go s.newClient(connection)
		}
	}
}