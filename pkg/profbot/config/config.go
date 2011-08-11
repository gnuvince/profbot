package config

import (
	"flag"
)

var Nickname string = "profbot"
var Password string
var Prefix string = "~"
var Channel string = "#chan"
var Server string = "localhost"
var Port string = "6667"
var DatabaseName string = "profbot.db"


func GetFlags() *flag.FlagSet {
	f := flag.NewFlagSet("progbot", flag.ExitOnError)
	f.StringVar(&Nickname, "nickname", Nickname, "nickname of the bot")
	f.StringVar(&Password, "password", Password, "nickserv password of the bot")
	f.StringVar(&Prefix, "prefix", Prefix, "prefix to trigger bot commands")
	f.StringVar(&Channel, "channel", Channel, "channel to go on")
	f.StringVar(&Server, "server", Server, "server to connect to")
	f.StringVar(&Port, "port", Port, "port of the server")
	f.StringVar(&DatabaseName, "database", DatabaseName, "database name")
	return f
}
