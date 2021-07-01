package main

import "net"

type room struct {
	roomName string
	members  map[net.Addr]*client
}

func (r *room) broadcast(c *client, msg string) {
	for addr, m := range r.members {
		if c.conn.RemoteAddr() != addr {
			m.clientMsg(msg)
		}
	}
}
