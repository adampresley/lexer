package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

/*
Lexer object contains the state of our parser and provides
a stream for accepting tokens.

Based on work by Rob Pike
http://cuddle.googlecode.com/hg/talk/lex.html#landing-slide
*/
type Lexer struct {
	Name   string
	Input  string
	Tokens chan Token
	State  LexFn

	Start int
	Pos   int
	Width int
}

/*
Backup puts the position tracker back to the beginning of the last read token.
*/
func (lexer *Lexer) Backup() {
	lexer.Pos -= lexer.Width
}

/*
CurrentCharacter returns the current character at the position tracker
*/
func (lexer *Lexer) CurrentCharacter() string {
	return lexer.Input[lexer.Pos : lexer.Pos+1]
}

/*
CurrentInput returns a slice of the current input from the current lexer start position
to the current position.
*/
func (lexer *Lexer) CurrentInput() string {
	return lexer.Input[lexer.Start:lexer.Pos]
}

/*
Dec dsecrement the position tracker back a single character
*/
func (lexer *Lexer) Dec() {
	lexer.Pos--
}

/*
Emit puts a token onto the token channel. The value of this token is
read from the input based on the current lexer position.
*/
func (lexer *Lexer) Emit(tokenType TokenType) {
	lexer.Tokens <- Token{Type: tokenType, Value: lexer.Input[lexer.Start:lexer.Pos]}
	lexer.Start = lexer.Pos
}

/*
EmitWithTransform allows you to put a typed-token onto the channel. The value
is read from the input based on the current lexer position, and then
passed to a provided transform function. That is then placed on the token
channel.
*/
func (lexer *Lexer) EmitWithTransform(tokenType TokenType, transformFn TokenValueTransformer) {
	lexer.Tokens <- Token{Type: tokenType, Value: transformFn(lexer.Input[lexer.Start:lexer.Pos])}
	lexer.Start = lexer.Pos
}

/*
Errorf returns a token with error information. This conforms to the
LexFn type
*/
func (lexer *Lexer) Errorf(format string, args ...interface{}) LexFn {
	lexer.Tokens <- Token{
		Type:  TOKEN_ERROR,
		Value: fmt.Sprintf(format, args...),
	}

	return nil
}

/*
Ignore disregards the current token by setting the lexer's start
position to the current reading position.
*/
func (lexer *Lexer) Ignore() {
	lexer.Start = lexer.Pos
}

/*
Inc move the position tracker forward one character
*/
func (lexer *Lexer) Inc() {
	lexer.Pos++

	if lexer.Pos > utf8.RuneCountInString(lexer.Input) {
		lexer.Pos--
	}
}

/*
InputToEnd returns a slice of the input from the current lexer position
to the end of the input string.
*/
func (lexer *Lexer) InputToEnd() string {
	return lexer.Input[lexer.Pos:]
}

/*
IsEOF returns true if the lexer is at the end of the
input stream.
*/
func (lexer *Lexer) IsEOF() bool {
	return lexer.Pos >= utf8.RuneCountInString(lexer.Input)
}

/*
IsNewline returns true if the current character is a newline character
*/
func (lexer *Lexer) IsNewline() bool {
	return lexer.CurrentCharacter() == "\n"
}

/*
IsNumber returns true if the current character is a number
*/
func (lexer *Lexer) IsNumber() bool {
	ch, _ := utf8.DecodeRuneInString(lexer.Input[lexer.Pos:])
	return unicode.IsNumber(ch)
}

/*
IsWhitespace returns true if then current character is whitespace
*/
func (lexer *Lexer) IsWhitespace() bool {
	ch, _ := utf8.DecodeRuneInString(lexer.Input[lexer.Pos:])
	return unicode.IsSpace(ch)
}

/*
Next reads the next rune (character) from the input stream
and advances the lexer position.
*/
func (lexer *Lexer) Next() rune {
	if lexer.Pos >= utf8.RuneCountInString(lexer.Input) {
		lexer.Width = 0
		return EOF
	}

	result, width := utf8.DecodeRuneInString(lexer.Input[lexer.Pos:])

	lexer.Width = width
	lexer.Pos += lexer.Width
	return result
}

/*
NextToken returns the next token from the channel
*/
func (lexer *Lexer) NextToken() Token {
	return <-lexer.Tokens
}

/*
Peek returns the next rune in the stream, then puts the lexer
position back. Basically reads the next rune without consuming
it.
*/
func (lexer *Lexer) Peek() rune {
	rune := lexer.Next()
	lexer.Backup()
	return rune
}

/*
PeekCharacters returns what the next set of characters in the input
stream is.
*/
func (lexer *Lexer) PeekCharacters(numCharacters int) string {
	end := lexer.Pos + numCharacters
	if end > utf8.RuneCountInString(lexer.Input) {
		end = utf8.RuneCountInString(lexer.Input)
	}

	return lexer.Input[lexer.Pos:end]
}

/*
Run starts the lexical analysis and feeding tokens into the
token channel.
*/
func (lexer *Lexer) Run() {
	go func() {
		for {
			lexer.State = lexer.State(lexer)
			if lexer.State == nil {
				break
			}
		}

		lexer.Shutdown()
	}()
}

/*
Shutdown closes up the token stream
*/
func (lexer *Lexer) Shutdown() {
	close(lexer.Tokens)
}

/*
SkipWhitespace skips whitespace characters until we get something meaningful.
*/
func (lexer *Lexer) SkipWhitespace() {
	var ch rune

	for {
		ch = lexer.Next()

		if !unicode.IsSpace(ch) {
			lexer.Dec()
			lexer.Start = lexer.Pos
			break
		}

		if ch == EOF {
			lexer.Emit(TOKEN_EOF)
			break
		}
	}
}
