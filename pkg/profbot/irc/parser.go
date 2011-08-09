package irc


func Parse(rawMessage string) (message Message, ok bool) {
	tokens, scanOk := Tokenize(rawMessage)
	if !scanOk {
		return
	}
	for _, token := range tokens {
		switch token.Type {
		case Nick:
			message.Nick = token.Text

		case Host:
			message.Host = token.Text

		case Command:
			message.Command = token.Text

		case Parameter:
			message.Parameters = append(message.Parameters, token.Text)

		case RestParameter:
			message.Parameters = append(message.Parameters, token.Text)
		}
	}

	return message, true
}
