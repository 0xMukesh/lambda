package parser

import (
	"errors"
	"fmt"

	"github.com/0xmukesh/lambda/internal/ast"
	"github.com/0xmukesh/lambda/internal/tokens"
)

type Parser struct {
	Tokens  []*tokens.Token
	Index   int
	context []string
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

func (p *Parser) pushVar(name string) {
	p.context = append(p.context, name)
}

func (p *Parser) popVar() {
	p.context = p.context[:len(p.context)-1]
}

func (p *Parser) resolve(name string) int {
	for i := len(p.context) - 1; i >= 0; i-- {
		if p.context[i] == name {
			return len(p.context) - i - 1
		}
	}

	return -1
}

// term ::= LAMBDA IDENTIFIER DOT term | application
func (p *Parser) term(node ast.Node) (ast.Node, error) {
	next := p.peek()
	if next == nil {
		return node, nil
	}

	if next.Type == tokens.Lambda {
		p.advance()

		varName, err := p.consume(tokens.Variable)
		if err != nil {
			return nil, err
		}

		if err := p.expect(tokens.Dot); err != nil {
			return nil, err
		}

		p.pushVar(string(varName.Char))
		term, err := p.term(node)
		p.popVar()

		if err != nil {
			return nil, err
		}
		if term == nil {
			return nil, errors.New("missing lambda body")
		}

		return &ast.Abstraction{
			Param: string(varName.Char),
			Body:  term,
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
		p.advance()
		term, err := p.term(node)
		if err != nil {
			return nil, err
		}

		if err := p.expect(tokens.RightParen); err != nil {
			return nil, err
		}

		return term, nil
	case tokens.Variable:
		p.advance()
		idx := p.resolve(string(next.Char))

		return &ast.Variable{
			Value: next.Char,
			Index: idx,
		}, nil
	case tokens.NamedTerm:
		p.advance()
		return &ast.NamedTermRef{
			Name: next.Value,
		}, nil
	default:
		return nil, nil
	}
}

func (p *Parser) Parse() (ast.Node, error) {
	if p.Index >= len(p.Tokens) {
		return nil, nil
	}

	if len(p.Tokens) > 2 && p.Tokens[0].Type == tokens.NamedTerm && p.Tokens[1].Type == tokens.Assign {
		name := p.Tokens[0].Value
		p.advance()
		p.advance()

		body, err := p.term(nil)
		if err != nil {
			return nil, err
		}
		if body == nil {
			return nil, errors.New("missing rhs side for assignment")
		}

		return &ast.Assignment{
			Name: name,
			Body: body,
		}, nil
	}

	return p.term(nil)
}
