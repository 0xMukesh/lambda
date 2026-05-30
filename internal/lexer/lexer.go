package lexer

import (
	"fmt"
	"unicode/utf8"

	"github.com/0xmukesh/lambda/internal/tokens"
)

type Lexer struct {
	Source []byte
	Index  int
	Char   byte
}

func NewLexer(source []byte) *Lexer {
	return &Lexer{
		Source: source,
		Index:  0,
		Char:   0,
	}
}

func (l *Lexer) Lex() (*tokens.Token, error) {
	var tok *tokens.Token
	if l.Index >= len(l.Source) {
		return tokens.NewToken(tokens.EOF, 0), nil
	}

	l.Char = l.Source[l.Index]
	l.Index++

	isValidUtf8 := utf8.ValidRune(rune(l.Char))
	if !isValidUtf8 {
		return nil, fmt.Errorf("got non-utf8 character: %c", l.Char)
	}

	if l.Char == '\n' || l.Char == ' ' || l.Char == '\t' || l.Char == 0 {
		return tok, nil
	}

	switch l.Char {
	case '\\':
		tok = tokens.NewToken(tokens.Lambda, l.Char)
	case '.':
		tok = tokens.NewToken(tokens.Dot, l.Char)
	case '(':
		tok = tokens.NewToken(tokens.LeftParen, l.Char)
	case ')':
		tok = tokens.NewToken(tokens.RightParen, l.Char)
	default:
		tok = tokens.NewToken(tokens.Identifier, l.Char)
	}

	return tok, nil
}

func (l *Lexer) LexAll() ([]*tokens.Token, error) {
	toks := []*tokens.Token{}

	for {
		tok, err := l.Lex()
		if err != nil {
			return nil, err
		}

		if tok == nil {
			continue
		}

		if tok.Type == tokens.EOF {
			return toks, nil
		}

		toks = append(toks, tok)
	}
}
