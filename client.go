package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type client struct {
	name     string
	conn     net.Conn
	room     *room
	commands chan<- command
}

func (c *client) readInput() {
	for {
		message, err := bufio.NewReader(c.conn).ReadString('\n')
		fmt.Println(message)
		if err != nil {
			log.Fatal("there is an error reading message", err)
		}
		message = strings.Trim(message, "\r\n")
		args := strings.Split(message, " ")
		fmt.Println(args)
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/name":
			c.commands <- command{
				id:     CMD_NAME,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- command{
				id:     CMD_JOIN,
				client: c,
				args:   args,
			}
		case "/rooms":
			c.commands <- command{
				id:     CMD_ROOMS,
				client: c,
			}
		case "/msg":
			c.commands <- command{
				id:     CMD_MSG,
				client: c,
				args:   args,
			}
		case "/quit":
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
			}
		default:
			c.clientErr(fmt.Errorf("unknown command: %s", cmd))
		}
	}
}

func (c *client) clientErr(err error) {
	c.conn.Write([]byte("error" + err.Error() + "\n"))
}
func (c *client) clientMsg(msg string) {
	c.conn.Write([]byte(">>> " + msg + "\n"))
}
