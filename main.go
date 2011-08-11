package main

import (
	"log"
	"os"
	"os/signal"

	"profbot/irc"
	"profbot/config"

	"profbot/plugins/pong"
	"profbot/plugins/nickserv"
	"profbot/plugins/url"
	"profbot/plugins/seen"
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

	c.Register(pong.New())
	c.Register(nickserv.New())
	c.Register(url.New(":memory:"))
	c.Register(seen.New(":memory:"))

	c.Nick(config.Nickname)
	c.User(config.Nickname, config.Nickname)
	c.Join(config.Channel)

	c.Loop()
}
