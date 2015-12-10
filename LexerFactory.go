package lexer

/*
NewLexer starts a new lexer with a given input string. This returns the
instance of the lexer and a channel of tokens. Reading this stream
is the way to parse a given input and perform processing.
*/
func NewLexer(name string, input string, startFn LexFn) *Lexer {
	l := &Lexer{
		Name:   name,
		Input:  input,
		State:  startFn,
		Tokens: make(chan Token, 100),
	}

	return l
}
