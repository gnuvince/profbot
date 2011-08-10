package irc

type UserCommand struct {
	Name string
	Help string
	Fn func(*Connection, *Message)
}
