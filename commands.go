package main

type commandID int

const (
	comm_name commandID = iota
	comm_join
	comm_rooms
	comm_message
	comm_quit
	comm_help
	comm_lvroom
	comm_members
)

type command struct {
	id        commandID
	client    *client
	arguments []string
}

const HELP_DISPLAY_STRING = `
-help : display help
-nick [ new nick ]: change nick
-rooms : check for open rooms
-join [ room name ] : join/create pointed room
-lvroom : leave current room
-members : display members in current room
-quit : disconnect
`