package irc

type UserCommand func(*Connection, *Message)
