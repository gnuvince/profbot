package main

import (
	"log"
	"os"
	"os/signal"

	"profbot/irc"
	"profbot/config"
	"profbot/db"

	"profbot/plugins/pong"
	"profbot/plugins/nickserv"
	"profbot/plugins/url"
	"profbot/plugins/seen"
)


func main() {
	flags := config.GetFlags()
	flags.Parse(os.Args[1:])

	c, err := irc.Connect(config.Server, config.Port)
	defer c.Close()


	if err != nil {
		log.Fatalf("Cannot establish IRC connection")
	}


	go func() {
		for sig := range signal.Incoming {
			if sig == os.SIGINT {
				c.Shutdown()
			}
		}
	}()

	dbConn := db.Open(config.DatabaseName)
	defer dbConn.Close()

	c.Register(pong.New())
	c.Register(nickserv.New())
	c.Register(url.New(dbConn))
	c.Register(seen.New(dbConn))

	c.Nick(config.Nickname)
	c.User(config.Nickname, config.Nickname)
	c.Join(config.Channel)

	c.Loop()
}
