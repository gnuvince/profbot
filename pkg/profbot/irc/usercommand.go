package irc


type UserCommand struct {
	Name string
	Fn func(*Connection, *Message)
}


func (c *Connection) RegisterCommand(name string, u UserCommand) {
	c.commands[name] = u
}
