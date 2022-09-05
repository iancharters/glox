package scanner

import (
	"strconv"
	"unicode"

	"github.com/iancharters/glox/errors"
	"github.com/iancharters/glox/token"
)

const NULL_TERMINATOR = '\x00'

type Scanner struct {
	source  string
	tokens  []token.Token
	start   int
	current int
	line    int
}

func New(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  []token.Token{},
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() ([]token.Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		if err := s.scanToken(); err != nil {
			return nil, err
		}
	}

	s.tokens = append(s.tokens, token.Token{token.EOF, "", nil, s.line})

	return s.tokens, nil
}

func (s *Scanner) scanToken() error {
	c := s.advance()

	switch c {
	case '(':
		s.addToken(token.LEFT_PAREN)
	case ')':
		s.addToken(token.RIGHT_PAREN)
	case '{':
		s.addToken(token.LEFT_BRACE)
	case '}':
		s.addToken(token.RIGHT_BRACE)
	case ',':
		s.addToken(token.COMMA)
	case '.':
		s.addToken(token.DOT)
	case '-':
		s.addToken(token.MINUS)
	case '+':
		s.addToken(token.PLUS)
	case ';':
		s.addToken(token.SEMICOLON)
	case '*':
		s.addToken(token.STAR)
	case '!':
		if s.match('=') {
			s.addToken(token.BANG_EQUAL)
		} else {
			s.addToken(token.BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(token.EQUAL_EQUAL)
		} else {
			s.addToken(token.EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(token.LESS_EQUAL)
		} else {
			s.addToken(token.LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(token.GREATER_EQUAL)
		} else {
			s.addToken(token.GREATER)
		}
	case '/':
		if s.match('/') {
			// a comment that goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			for s.peek() != '*' && s.peekNext() != '/' {
				if s.isAtEnd() {
					return errors.New(s.line, "", errors.ErrUnterminatedComment)
				}

				s.advance()
			}
		} else {
			s.addToken(token.SLASH)
		}
	case ' ':
		break
	case '\r':
		break
	case '\t':
		break
	case '\n':
		s.line += 1
	case '"':
		if err := s.string(); err != nil {
			return err
		}
	case 'o':
		if s.match('r') {
			s.addToken(token.OR)
		}

	default:
		if isDigit(c) {
			if err := s.number(); err != nil {
				return err
			}
		} else if isAlpha(c) {
			s.identifier()
		} else {
			return errors.New(s.line, "", errors.ErrUnexpectedCharacter)
		}
	}

	return nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() rune {
	r := rune(s.source[s.current])
	s.current += 1
	return r
}

func (s *Scanner) addToken(t token.TokenType) {
	s.addTokenLiteral(t, nil)
}

func (s *Scanner) addTokenLiteral(t token.TokenType, l token.Literal) {
	lexeme := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.Token{t, lexeme, l, s.line})
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if rune(s.source[s.current]) != expected {
		return false
	}

	s.current += 1

	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return NULL_TERMINATOR
	}

	return rune(s.source[s.current])
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return NULL_TERMINATOR
	}

	return rune(s.source[s.current+1])
}

func (s *Scanner) string() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line += 1
		}

		s.advance()
	}

	if s.isAtEnd() {
		return errors.New(s.line, "", errors.ErrUnterminatedString)
	}

	s.advance() // the closing "

	literal := s.source[s.start+1 : s.current-1] // trim quotes
	s.addTokenLiteral(token.STRING, literal)

	return nil
}

func (s *Scanner) number() error {
	for isDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance() // consume the '.'

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	str := s.source[s.start:s.current]
	number, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return errors.New(s.line, "", errors.ErrParseFailure)
	}

	s.addTokenLiteral(token.NUMBER, number)
	return nil
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	lexeme := s.source[s.start:s.current]

	if kw, ok := token.Keywords[lexeme]; ok {
		s.addToken(kw)
	} else {
		s.addToken(token.IDENTIFIER)
	}
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func isAlpha(r rune) bool {
	return unicode.IsLetter(r)
}

func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || isDigit(r)
}
