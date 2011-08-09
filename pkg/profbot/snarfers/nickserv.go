package snarfers

import (
	"profbot/irc"
	"profbot/config"
)


func Nickserv(conn *irc.Connection, msg *irc.Message) {
	if msg.Command == "376" {
		conn.Privmsg("nickserv", "identify " + config.Password)
	}
}
