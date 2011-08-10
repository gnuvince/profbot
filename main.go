package main

import (
	"log"
	"os"
	"os/signal"

	"profbot/irc"
	"profbot/config"

	"profbot/plugins/pong"
	"profbot/plugins/nickserv"
)


func main() {
	flags := config.GetFlags()
	flags.Parse(os.Args[1:])

	c, err := irc.Connect("localhost", "6667")
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

	nickserv.Register(c)
	pong.Register(c)

	c.Nick(config.Nickname)
	c.User(config.Nickname, config.Nickname)
	c.Join(config.Channel)

	c.Loop()
}
