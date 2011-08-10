package irc

import (
	"bufio"
	"net"
	"testing"
)


func WithConn(f func(*Connection, *bufio.Reader)) {
	a, b := net.Pipe()
	buf := bufio.NewReader(b)
	conn := New(a)
	go conn.Loop()
	f(conn, buf)
}


func TestSendf(t *testing.T) {
	WithConn(func (c *Connection, res *bufio.Reader) {
		c.Sendf("A")
		if s, _ := res.ReadString('\n'); s != "A\r\n" {
			t.Errorf("c.Sendf() should add a \\r\\n")
		}
	})
}


func TestPrivmsg(t *testing.T) {
	WithConn(func (c *Connection, res *bufio.Reader) {
		c.Privmsg("#chan", "hello")
		expected := "PRIVMSG #chan :hello\r\n"
		if s, _ := res.ReadString('\n'); s != expected {
			t.Errorf("c.Privmsg(): expected %q, got %q", expected, s)
		}
	})
}


func TestNick(t *testing.T) {
	WithConn(func (c *Connection, res *bufio.Reader) {
		c.Nick("newnick")
		expected := "NICK newnick\r\n"
		if s, _ := res.ReadString('\n'); s != expected {
			t.Errorf("c.Nick(): expected %q, got %q", expected, s)
		}
	})
}


func TestUser(t *testing.T) {
	WithConn(func (c *Connection, res *bufio.Reader) {
		c.User("ident", "name")
		expected := "USER ident 12 * name\r\n"
		if s, _ := res.ReadString('\n'); s != expected {
			t.Errorf("c.User(): expected %q, got %q", expected, s)
		}
	})
}

func TestJoin(t *testing.T) {
	WithConn(func (c *Connection, res *bufio.Reader) {
		c.Join("#chan")
		expected := "JOIN #chan\r\n"
		if s, _ := res.ReadString('\n'); s != expected {
			t.Errorf("c.Join(): expected %q, got %q", expected, s)
		}
	})
}

func TestQuit(t *testing.T) {
	WithConn(func (c *Connection, res *bufio.Reader) {
		c.Quit("quitting")
		expected := "QUIT quitting\r\n"
		if s, _ := res.ReadString('\n'); s != expected {
			t.Errorf("c.Quit(): expected %q, got %q", expected, s)
		}
	})
}

func TestPong(t *testing.T) {
	WithConn(func (c *Connection, res *bufio.Reader) {
		c.Pong("server")
		expected := "PONG :server\r\n"
		if s, _ := res.ReadString('\n'); s != expected {
			t.Errorf("c.Pong(): expected %q, got %q", expected, s)
		}
	})
}


