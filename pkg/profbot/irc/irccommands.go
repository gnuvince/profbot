package irc

import (
	"fmt"
)



func (c *Connection) Sendf(s string, v ...interface{}) {
	c.out <- fmt.Sprintf(s, v...)
}


func (c *Connection) Privmsg(target, msg string) {
	c.Sendf("PRIVMSG %s :%s", target, msg)
}


func (c *Connection) Nick(nick string) {
	c.Sendf("NICK %s", nick)
}


func (c *Connection) User(ident, name string) {
	c.Sendf("USER %s 12 * %s", ident, name)
}


func (c *Connection) Join(channel string) {
	c.Sendf("JOIN %s", channel)
}


func (c *Connection) Quit(reason string) {
	c.Sendf("QUIT %s", reason)
}
