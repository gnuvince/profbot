package irc

import (
	"fmt"
	"regexp"
)

// An IRC message consists of:
// - The sender of the message (could be empty)
// - The host of the sender (could be empty)
// - The command sent
// - The list of parameters sent
type Message struct {
	Nick string
	Host string
	Command string
	Parameters []string
}


// The last parameter of an IRC message usually contains
// the text of the message (e.g.: the text typed by another
// client).  This is a convenience method to get that parameter.
// It returns nil if the list of parameters is empty, but
// in practice this should not happen.
func (m *Message) LastParameter() *string {
	if len(m.Parameters) != 0 {
		return &m.Parameters[len(m.Parameters) - 1]
	}
	return nil
}


// Return the user or channel to whom the message is destined.
func (m *Message) Target() *string {
    if len(m.Parameters) > 0 {
        return &m.Parameters[0]
    }
    return nil
}


// Break the last parameter (usually the text of a PRIVMSG) into
// words
var splitter *regexp.Regexp = regexp.MustCompile(`[^ \t]+`)
func (m *Message) Split() []string {
	parts := splitter.FindAllString(*m.LastParameter(), -1)
	return parts
}


func (m *Message) String() string {
	switch m.Command {
	case "JOIN":
		return fmt.Sprintf("* %s has joined %s", m.Nick, *m.Target())
	case "PART":
		return fmt.Sprintf("* %s has left %s", m.Nick, *m.Target())
	case "QUIT":
		return fmt.Sprintf("* %s has quit: %s", m.Nick, *m.LastParameter())
	case "NICK":
		return fmt.Sprintf("* %s is now known as %s", m.Nick, *m.LastParameter())
	case "TOPIC":
		return fmt.Sprintf("* %s has changed the topic of %s: %s",
			m.Nick, *m.Target(), *m.LastParameter())
	case "PRIVMSG":
		if m.Parameters[1] == "ACTION" {
			return fmt.Sprintf("* %s %s", m.Nick, *m.LastParameter())
		} else {
			return fmt.Sprintf("<%s> %s", m.Nick, *m.LastParameter())
		}
	}
	return fmt.Sprintf("<UNKNOWN MESSAGE TYPE: %s>", m.Command)
}
