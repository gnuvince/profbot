package snarfers

import (
	"profbot/irc"
)


func Pong(conn *irc.Connection, msg *irc.Message) {
	if msg.Command == "PING" {
		conn.Pong(*msg.LastParameter())
	}
}
