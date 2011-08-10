package nickserv

import (
	"profbot/irc"
	"profbot/config"
)


func Register(c *irc.Connection) {
	c.RegisterSnarfer(snarf)
}


func snarf(conn *irc.Connection, msg *irc.Message) {
	if msg.Command == "376" && config.Password != "" {
		conn.Privmsg("nickserv", "identify " + config.Password)
	}
}
