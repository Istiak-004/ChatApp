package main

const (
	CMD_NAME int = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
)

type command struct {
	id     int
	client *client
	args   []string
}
