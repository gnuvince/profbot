package url

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"time"
	sqlite "gosqlite.googlecode.com/hg/sqlite"

	"profbot/irc"
	"profbot/db"
)

const (
	tableName = "urls"
	maxResponses = 10
)


var urlPattern *regexp.Regexp = regexp.MustCompile(`(((ftp|https?)://[^ ]+)|(www\.[^ ]+))`)



type t struct {
	db *sqlite.Conn
}


func New(dbname string) *t {
	plugin := new(t)
	plugin.db = db.Open(dbname)

	if !db.TableExists(plugin.db, tableName) {
		log.Printf("Creating table %s", tableName)
		db.CreateTable(plugin.db, tableName,
			"id integer primary key",
			"url text",
			"nick text",
			"datetime integer",
			"rawmessage text")
	}

	return plugin
}


func (s *t) Snarfers() []irc.Snarfer {
	return []irc.Snarfer{
		func (c *irc.Connection, m *irc.Message) {
			s.snarf(c, m)
		},
	}
}


func (s *t) Commands() []irc.UserCommand {
	return []irc.UserCommand{
		irc.UserCommand{"url",
			func (c *irc.Connection, m *irc.Message) {
				s.lastUrl(c, m)
			},
		},
	}
}



func (s *t) snarf(conn *irc.Connection, msg *irc.Message) {
	if msg.Command != "PRIVMSG" {
		return
	}

	matches := urlPattern.FindAllString(*msg.LastParameter(), -1)
	stmt, err := s.db.Prepare(
		fmt.Sprintf("insert into %s (url, nick, datetime, rawmessage) values(?, ?, ?, ?)", tableName))
	defer stmt.Finalize()

	if err != nil {
		return
	}

	for _, match := range matches {
		stmt.Exec(match, msg.Nick, time.Seconds(), *msg.LastParameter())
		stmt.Next()
		stmt.Reset()
	}
}



func (s *t) lastUrl(conn *irc.Connection, msg *irc.Message) {
	if msg.Command != "PRIVMSG" {
		return
	}

	var nFlag int
	var contextFlag string
	var fromFlag string
	var contains string = ""

	fs := flag.NewFlagSet("urls", flag.ContinueOnError)
	fs.IntVar(&nFlag, "n", 1, "number of results")
	fs.StringVar(&fromFlag, "from", "", "sender's nick contains string")
	fs.StringVar(&contextFlag, "context", "", "context of the message")
	fs.Usage = func() {
		conn.Privmsg(*msg.Target(), "usage: url [-from=<nick>] [-n=<number to display>] [-context=<string>] [<substring>]")
	}

	parts := msg.Split()
	if fs.Parse(parts[1:]) != nil {
		return
	}

	if fs.NArg() > 0 {
		contains = fs.Arg(0)
	}

 	query := fmt.Sprintf(
		"SELECT url FROM %s WHERE nick LIKE ? AND url LIKE ? AND rawmessage LIKE ? ORDER BY datetime DESC LIMIT ?", tableName)
	stmt, err := s.db.Prepare(query)
	defer stmt.Finalize()
	if err != nil {
		conn.Privmsg(*msg.Target(), "sql error")
		return
	}

	stmt.Exec("%" + fromFlag + "%",
		"%" + contains + "%",
		"%" + contextFlag + "%",
		nFlag)
	more := stmt.Next()
	var output string

	if !more {
		output = "no URL matches those criterias"
	}

	left := maxResponses
	for more && left > 0 {
		var url string

		stmt.Scan(&url)
		output += url + "   "

		more = stmt.Next()
		left--
	}

	conn.Privmsg(*msg.Target(), output)

}
