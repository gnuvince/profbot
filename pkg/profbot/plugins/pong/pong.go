package pong

import (
	"profbot/irc"
)

func Register(c *irc.Connection) {
	c.RegisterSnarfer(snarf)
}

func snarf(conn *irc.Connection, msg *irc.Message) {
	if msg.Command == "PING" {
		conn.Pong(*msg.LastParameter())
	}
}
