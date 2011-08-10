package pong

import (
	"testing"
	"net"
	"bufio"

	"profbot/irc"
)


func TestSnarfMessage(t *testing.T) {
	in, server := net.Pipe()
	bufin := bufio.NewReader(in)
	conn := irc.New(server)
	defer conn.Shutdown()
	go conn.Loop()

	snarf(conn, &irc.Message{"", "", "PING", []string{"server"}})
	s, _ := bufin.ReadString('\n')

	expected := "PONG :server\r\n"
	if s != expected {
		t.Errorf("Expected %q, got %q", expected, s)
	}
}
