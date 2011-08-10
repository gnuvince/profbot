package pong

import (
	"profbot/irc"
)

type t struct {}

func New() *t {
	return new(t)
}

func (u *t) Snarfers() []irc.Snarfer {
	return []irc.Snarfer{snarf}
}

func (u *t) Commands() []irc.UserCommand {
	return []irc.UserCommand{}
}

func snarf(conn *irc.Connection, msg *irc.Message) {
	if msg.Command == "PING" {
		conn.Pong(*msg.LastParameter())
	}
}
