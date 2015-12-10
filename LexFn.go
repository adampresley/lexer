package lexer

/*
LexFn defines a function type that lexer parsing functions must implement. These
functions are what parse input text
*/
type LexFn func(*Lexer) LexFn
