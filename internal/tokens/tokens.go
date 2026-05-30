package tokens

import "fmt"

type TokenType int

const (
	Lambda TokenType = iota
	Dot
	LeftParen
	RightParen
	Identifier
	EOF
)

func (tt TokenType) String() string {
	switch tt {
	case Lambda:
		return "LAMBDA"
	case Dot:
		return "DOT"
	case LeftParen:
		return "LEFT_PAREN"
	case RightParen:
		return "RIGHT_PAREN"
	case Identifier:
		return "IDENTIFIER"
	case EOF:
		return "EOF"
	default:
		panic(fmt.Sprintf("invalid token type: %d", tt))
	}
}

type Token struct {
	Type TokenType
	Char byte
}

func NewToken(tokenType TokenType, char byte) *Token {
	return &Token{
		Type: tokenType,
		Char: char,
	}
}
