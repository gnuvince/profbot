package main

import (
	"log"
	"os"
	"os/signal"

	"profbot/irc"
)


func main() {
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


	c.Nick("profbot")
	c.User("profbot", "profbot")
	c.Join("#programmeur")

	c.Loop()
}
