package irc

// States
const (
	Start int = iota
	GetNick
	GetHost
	GetCommand
	GetParams
	GetParam
	GetRest
)


// Token types
const (
	Nick = iota
	Host
	Command
	Parameter
	RestParameter
)

func typeToString(t int) string {
	return []string{
		"Nick", "Host", "Command", "Parameter", "RestParameter",
	}[t]
}

type Token struct {
	Type int
	Text string
}

func TokensEqual(t1, t2 Token) bool {
	return t1.Type == t2.Type &&
		t1.Text == t2.Text
}


func isAlpha(c int) bool {
	return c >= 'A' && c <= 'Z' ||
		c >= 'a' && c <= 'z'
}

func isDigit(c int) bool {
	return c >= '0' && c <= '9'
}

func isAlnum(c int) bool {
	return isAlpha(c) || isDigit(c)
}


func Tokenize(input string) (tokens []Token, ok bool) {
	addToken := func(token Token) {
		tokens = append(tokens, token)
	}

	token := ""
	state := Start

	for _, c := range input {
		switch state {
		case Start:
			if c == ':' {
				state = GetNick
			} else if isAlpha(c) {
				state = GetCommand
				token += string(c) // We've already consumed the first character of the command
			} else {
				return tokens, false
			}

		case GetNick:
			if c != ' ' && c != '!' {
				token += string(c)
			} else if c == '!' {
				addToken(Token{Nick, token})
				token = ""
				state = GetHost
			} else if c == ' ' {
				addToken(Token{Nick, token})
				token = ""
				state = GetCommand
			}

		case GetHost:
			if c != ' ' {
				token += string(c)
			} else {
				addToken(Token{Host, token})
				token = ""
				state = GetCommand
			}

		case GetCommand:
			if isAlnum(c) {
				token += string(c)
			} else if c == ' ' {
				addToken(Token{Command, token})
				token = ""
				state = GetParams
			} else {
				return tokens, false
			}

		case GetParams:
			if c != ':' {
				state = GetParam
				token += string(c) // We've already consumed the first character of the command
			} else {
				state = GetRest
			}

		case GetParam:
			if c != ' ' {
				token += string(c)
			} else {
				addToken(Token{Parameter, token})
				token = ""
				state = GetParams
			}

		case GetRest:
			if c != '\r' && c != '\n' {
				token += string(c)
			} else {
				ok = true
			}
		}
	}
	addToken(Token{RestParameter, token})
	return tokens, ok
}
