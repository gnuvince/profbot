package irc

type Snarfer func(*Connection, *Message)


func (c *Connection) RegisterSnarfer(sn Snarfer) {
	c.snarfers = append(c.snarfers, sn)
}
