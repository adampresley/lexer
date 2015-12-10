package lexer

/*
A TokenType defines the types of tokens available. Create your own to describe
your input
*/
type TokenType int

const (
	TOKEN_ERROR TokenType = -2
	TOKEN_EOF   TokenType = -1
)
