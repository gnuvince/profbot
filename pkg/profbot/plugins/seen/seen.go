package seen

import (
	"fmt"
	"flag"
	"log"
	"time"
	sqlite "gosqlite.googlecode.com/hg/sqlite"

	"profbot/irc"
	"profbot/db"
)


const (
	tableName = "seen"
)


type t struct {
	db *sqlite.Conn
}


func New(dbConn *sqlite.Conn) *t {
	plugin := new(t)
	plugin.db = dbConn

	if !db.TableExists(plugin.db, tableName) {
		log.Printf("Creating table %s", tableName)
		db.CreateTable(plugin.db, tableName,
			"id integer primary key",
			"datetime integer",
			"nick text",
			"message text")
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
		irc.UserCommand{"seen",
			func (c *irc.Connection, m *irc.Message) {
				s.seen(c, m)
			},
		},
	}
}



func (s *t) snarf(conn *irc.Connection, msg *irc.Message) {
	stmt, err := s.db.Prepare(
		fmt.Sprintf("select id from %s where nick like ?", tableName))
	defer stmt.Finalize()

	if err != nil {
		log.Printf("seen.snarf() 1: sql error: %s", err)
		return
	}

	stmt.Exec(msg.Nick)
	found := stmt.Next()

	if found {
		var id int
		stmt.Scan(&id)
		stmt.Reset()
		stmt, err = s.db.Prepare(
			fmt.Sprintf("update %s set datetime = ?, message = ? where id = ?", tableName))

		if err != nil {
			log.Printf("seen.snarf() 2: sql error: %s", err)
			return
		}

		stmt.Exec(time.Seconds(), msg.String(), id)
		stmt.Next()
	} else {
		stmt.Reset()
		stmt, err = s.db.Prepare(
			fmt.Sprintf("insert into %s (datetime, nick, message) values(?, ?, ?)", tableName))

		if err != nil {
			log.Printf("seen.snarf() 3: sql error: %s", err)
			return
		}

		stmt.Exec(time.Seconds(), msg.Nick, msg.String())
		stmt.Next()
	}
}


func (s *t) seen(conn *irc.Connection, msg *irc.Message) {
	if msg.Command != "PRIVMSG" {
		return
	}

	fs := flag.NewFlagSet("seen", flag.ContinueOnError)
	fs.Usage = func() {
		conn.Privmsg(*msg.Target(), "usage: seen <nick>")
	}

	parts := msg.Split()
	if fs.Parse(parts[1:]) != nil {
		return
	}

	if fs.NArg() == 0 {
		fs.Usage()
		return
	}

	target := fs.Arg(0)
	stmt, err := s.db.Prepare(
		fmt.Sprintf("SELECT datetime, message FROM %s WHERE nick LIKE ? ORDER BY datetime desc LIMIT 1", tableName))
	defer stmt.Finalize()
	if err != nil {
		log.Printf("seen.seen(): sql error: %s", err)
		return
	}

	stmt.Exec(target)
	found := stmt.Next()
	if found {
		var datetime int64
		var rawmessage string
		stmt.Scan(&datetime, &rawmessage)
		t := time.SecondsToLocalTime(datetime)
        str := t.Format(time.UnixDate)
		conn.Privmsg(*msg.Target(), str + ": " + rawmessage)
	} else {
		conn.Privmsg(*msg.Target(), "I haven't seen " + target)
	}
}
