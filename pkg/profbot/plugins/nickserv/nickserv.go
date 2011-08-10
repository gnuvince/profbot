package nickserv

import (
	"profbot/irc"
	"profbot/config"
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
	if msg.Command == "376" && config.Password != "" {
		conn.Privmsg("nickserv", "identify " + config.Password)
	}
}
