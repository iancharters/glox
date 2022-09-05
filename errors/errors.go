package errors

import (
	"errors"
	"fmt"

	"github.com/iancharters/glox/token"
)

var (
	ErrUnexpectedCharacter = errors.New("unexpected character")
	ErrUnterminatedString  = errors.New("unterminated string")
	ErrUnterminatedComment = errors.New("unterminated multiline comment")

	ErrParseFailure = errors.New("failed to parse float")
)

type LoxError struct {
	Line  int
	Where string

	Err error
}

func New(line int, where string, err error) LoxError {
	return LoxError{line, where, err}
}

func (e LoxError) Error() string {
	return fmt.Sprintf("[line %d] Error %s: %s", e.Line, e.Where, e.Err.Error())
}

func NewUnexpectedCharacterError(line int) LoxError {
	return New(line, "", ErrUnexpectedCharacter)
}

func NewUnterminatedStringError(line int) LoxError {
	return New(line, "", ErrUnterminatedString)
}

func NewParseError(t token.Token, message string) LoxError {
	if t.Type == token.EOF {
		return New(t.Line, "at end", errors.New(message))
	}

	return New(t.Line, fmt.Sprintf("at '%s'", t.Lexeme), errors.New(message))
}
