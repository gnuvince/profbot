package irc


import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"profbot/config"
)

/*
 * A connection to an IRC server has the following properties:
 * - A socket connected to the IRC server (can be nil for testing)
 * - An buf buffer to communicate with the IRC server
 * - An input channel: after the raw strings from the inbound stream
 *   has been parsed, that Message is send on the inbound channel
 *   where it'll be dispatched.
 * - An ouput channel: the only way to send messages to the IRC server
 *   is by sending a string on this channel.
 * - A quit channel: once this channel receives an input, it shuts
     down the IRC connection.
 */
type Connection struct {
	sock net.Conn
	buf *bufio.ReadWriter
	in chan Message
	out chan string
	quit chan bool

	snarfers []Snarfer
	commands map[string]UserCommand
}



/*
 * Create a new Connection object.
 */
func New(c net.Conn) *Connection {
	conn := new(Connection)
	conn.in = make(chan Message, 64)
	conn.out = make(chan string, 64)
	conn.quit = make(chan bool)
	conn.snarfers = []Snarfer{}
	conn.commands = map[string]UserCommand{}
	conn.sock = c
	conn.buf = bufio.NewReadWriter(
		bufio.NewReader(conn.sock),
		bufio.NewWriter(conn.sock))

	return conn
}


/*
 * Connect to an IRC network by providing its hostname and port.
 * Connect() uses NewCustom internally.
 */
func Connect(hostname, port string) (*Connection, os.Error) {
	conn, err := net.Dial("tcp", hostname + ":" + port)
	if err != nil {
		return nil, err
	}
	return New(conn), err
}


/*
 * When a connection is closed, we close the socket as well as the
 * communication channels.
 */
func (c *Connection) Close() {
	log.Println("Closing connection...")

	c.sock.Close()

	close(c.in)
	close(c.out)
	close(c.quit)
}


/*
 * Send the shutdown signal to the connection.
 */
func (c *Connection) Shutdown() {
	c.Quit("leaving")
	c.quit <- true
}


/*
 * Launch the goroutines for receiving and sending data,
 * fetch messages on the inbound channel and dispatch.
 * If a message is received on the quit channel, exit the
 * loop.
 */
func (c *Connection) Loop() {
	go c.recv()
	go c.send()

	for {
		select {
		case m := <- c.in:
			fmt.Printf("%+v\n", m)

			if strings.HasPrefix(*m.LastParameter(), config.Prefix) {
				parts := m.Split()
				cmd, ok := c.commands[parts[0]]
				if ok {
					cmd.Fn(c, &m)
				}
			}

			for _, snarfer := range c.snarfers {
				go snarfer(c, &m)
			}

		case <- c.quit:
			return
		}
	}
}

/*
 * Read a string from the input buffer, parse it and send
 * the parsed message on the input channel.
 */
func (c *Connection) recv() {
	for {
		s, err := c.buf.ReadString('\n')
		if err != nil {
			log.Fatalf("c.recv(): %s", err.String())
		}

		if m, ok := Parse(s); ok {
			c.in <- m
		} else {
			log.Println("Parse error on: " + s)
		}
	}
}


/*
 * Take strings from the output channel and write them to the output
 * buffer.  Flush the output to make sure it always goes through.
 */
func (c *Connection) send() {
	for s := range c.out {
		s = strings.TrimSpace(s)
		if _, err := c.buf.WriteString(s + "\r\n"); err != nil {
			log.Fatalf("c.send(): %s", err.String())
		}
		c.buf.Flush()
	}
}




func (c *Connection) Register(p Plugin) {
	for _, snarfer := range p.Snarfers() {
		c.RegisterSnarfer(snarfer)
	}

	for _, command := range p.Commands() {
		c.RegisterCommand(config.Prefix + command.Name, command)
	}
}
