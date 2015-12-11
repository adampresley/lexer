package lexer

/*
A Token represents a parsed item in a source input. A token has a type
and a value. These are used to determine what to do next.
*/
type Token struct {
	Type  TokenType
	Value interface{}
}

func (token Token) IsEmpty() bool {
	return token.Type == 0 && token.Value == nil
}

func (token Token) String() string {
	switch token.Type {
	case TOKEN_EOF:
		return "EOF"

	case TOKEN_ERROR:
		return (token.Value).(string)
	}

	return (token.Value).(string)
}
