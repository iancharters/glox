package parser

import (
	"github.com/iancharters/glox/errors"
	"github.com/iancharters/glox/token"
)

type Parser struct {
	current int
	tokens  []token.Token
}

func New(tokens []token.Token) *Parser {
	return &Parser{0, tokens}
}

func (p *Parser) Expression() (Expr, error) {
	return p.Equality()
}

func (p *Parser) Equality() (Expr, error) {
	expression, err := p.Comparison()
	if err != nil {
		return nil, err
	}

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()

		right, err := p.Comparison()
		if err != nil {
			return nil, err
		}

		expression = NewBinaryExpr(expression, &operator, right)
	}

	return expression, nil
}

func (p *Parser) Comparison() (Expr, error) {
	expression, err := p.Term()
	if err != nil {
		return nil, err
	}

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()

		right, err := p.Term()
		if err != nil {
			return nil, err
		}

		expression = NewBinaryExpr(expression, &operator, right)
	}

	return expression, nil
}

func (p *Parser) Term() (Expr, error) {
	expression, err := p.Factor()
	if err != nil {
		return nil, err
	}

	for p.match(token.PLUS, token.MINUS) {
		operator := p.previous()

		right, err := p.Factor()
		if err != nil {
			return nil, err
		}

		expression = NewBinaryExpr(expression, &operator, right)
	}

	return expression, nil
}

func (p *Parser) Factor() (Expr, error) {
	expression, err := p.Unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.STAR, token.SLASH) {
		operator := p.previous()

		right, err := p.Unary()
		if err != nil {
			return nil, err
		}
		expression = NewBinaryExpr(expression, &operator, right)
	}

	return expression, nil
}

func (p *Parser) Unary() (Expr, error) {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()

		right, err := p.Unary()
		if err != nil {
			return nil, err
		}

		return NewUnaryExpr(&operator, right), nil
	}

	return p.Primary()
}

func (p *Parser) Primary() (Expr, error) {
	if p.match(token.FALSE) {
		return NewLiteralExpr(false), nil
	}

	if p.match(token.TRUE) {
		return NewLiteralExpr(true), nil
	}

	if p.match(token.NIL) {
		return NewLiteralExpr(nil), nil
	}

	if p.match(token.NUMBER, token.STRING) {
		return NewLiteralExpr(p.previous().Literal), nil
	}

	if p.match(token.LEFT_PAREN) {
		if err := p.consume(token.RIGHT_PAREN, "Expect ')' after expression."); err != nil {
			return nil, err
		}

		expression, err := p.Expression()
		if err != nil {
			return nil, err
		}

		return NewGroupingExpr(expression), nil
	}

	return nil, errors.NewParseError(p.peek(), "Expect expression.")
}

func (p *Parser) match(types ...token.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(t token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == t
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *Parser) consume(t token.TokenType, message string) error {
	if p.check(t) {
		p.advance()
		return nil
	}

	return errors.NewParseError(p.peek(), message)
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().Type == token.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case token.CLASS:
			fallthrough
		case token.FOR:
			fallthrough
		case token.FUN:
			fallthrough
		case token.IF:
			fallthrough
		case token.PRINT:
			fallthrough
		case token.RETURN:
			fallthrough
		case token.VAR:
			fallthrough
		case token.WHILE:
			return
		}

		p.advance()
	}
}

func (p *Parser) Parse() (Expr, error) {
	expression, err := p.Expression()
	if err != nil {
		return nil, err
	}

	return expression, nil
}
