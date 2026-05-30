package parser

import (
	"fmt"

	"github.com/0xmukesh/lambda/internal/ast"
	"github.com/0xmukesh/lambda/internal/tokens"
)

type Parser struct {
	Tokens []*tokens.Token
	Index  int
}

func NewParser(tokens []*tokens.Token) *Parser {
	return &Parser{
		Tokens: tokens,
		Index:  0,
	}
}

func (p *Parser) advance() *tokens.Token {
	if p.Index >= len(p.Tokens) {
		return nil
	}

	tok := p.Tokens[p.Index]
	p.Index++
	return tok
}

func (p *Parser) peek() *tokens.Token {
	if p.Index >= len(p.Tokens) {
		return nil
	}

	return p.Tokens[p.Index]
}

func (p *Parser) expect(tt tokens.TokenType) error {
	peek := p.peek()
	if peek == nil {
		return fmt.Errorf("expected %s, got EOF", tt)
	}
	if peek.Type != tt {
		return fmt.Errorf("expected %s, got %s", tt, peek.Type)
	}

	p.advance()
	return nil
}

func (p *Parser) consume(tt tokens.TokenType) (*tokens.Token, error) {
	peek := p.peek()
	if peek == nil {
		return nil, fmt.Errorf("expected %s, got EOF", tt)
	}
	if peek.Type != tt {
		return nil, fmt.Errorf("expected %s, got %s", tt, peek.Type)
	}

	return p.advance(), nil
}

// term ::= LAMBDA IDENTIFIER DOT term | application
func (p *Parser) term(node ast.Node) (ast.Node, error) {
	next := p.peek()
	if next == nil {
		return node, nil
	}

	if next.Type == tokens.Lambda {
		p.advance() // consuming lambda token

		ident, err := p.consume(tokens.Identifier)
		if err != nil {
			return nil, err
		}

		if err := p.expect(tokens.Dot); err != nil {
			return nil, err
		}

		term, err := p.term(node)
		if err != nil {
			return nil, err
		}

		return &ast.Abstraction{
			Param: &ast.Identifier{
				Value: string(ident.Char),
			},
			Body: term,
		}, nil
	}

	return p.application(node)
}

// application ::= atom application'
// application' ::= atom application' | EMPTY
func (p *Parser) application(node ast.Node) (ast.Node, error) {
	lhs, err := p.atom(node)
	if err != nil {
		return nil, err
	}

	for {
		rhs, err := p.atom(node)
		if err != nil {
			return nil, err
		}

		if rhs == nil {
			return lhs, nil
		} else {
			lhs = &ast.Application{
				Lhs: lhs,
				Rhs: rhs,
			}
		}
	}
}

// atom ::= LPAREN term RPAREN | IDENTIFIER
func (p *Parser) atom(node ast.Node) (ast.Node, error) {
	next := p.peek()
	if next == nil {
		return node, nil
	}

	switch next.Type {
	case tokens.LeftParen:
		p.advance() // consume left paren

		term, err := p.term(node)
		if err != nil {
			return nil, err
		}

		if err := p.expect(tokens.RightParen); err != nil {
			return nil, err
		}

		return term, nil
	case tokens.Identifier:
		p.advance() // consume identifier

		return &ast.Identifier{
			Value: string(next.Char),
		}, nil
	default:
		return nil, nil
	}
}

func (p *Parser) Parse() (ast.Node, error) {
	if p.Index >= len(p.Tokens) {
		return nil, nil
	}

	var node ast.Node
	result, err := p.term(node)
	if err != nil {
		return nil, err
	}

	return result, nil
}
