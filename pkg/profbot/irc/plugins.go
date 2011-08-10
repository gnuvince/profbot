package irc


type Plugin interface {
	Snarfers() []Snarfer
	Commands() []UserCommand
}
