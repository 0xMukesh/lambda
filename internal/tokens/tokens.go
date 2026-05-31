package tokens

import "fmt"

type TokenType int

const (
	Lambda TokenType = iota
	Variable
	NamedTerm
	Dot
	Assign
	LeftParen
	RightParen
	EOF
)

func (tt TokenType) String() string {
	switch tt {
	case Lambda:
		return "LAMBDA"
	case Variable:
		return "VARIABLE"
	case NamedTerm:
		return "NAMED_TERM"
	case Dot:
		return "DOT"
	case Assign:
		return "ASSIGN"
	case LeftParen:
		return "LEFT_PAREN"
	case RightParen:
		return "RIGHT_PAREN"
	case EOF:
		return "EOF"
	default:
		panic(fmt.Sprintf("invalid token type: %d", tt))
	}
}

type Token struct {
	Type  TokenType
	Char  byte
	Value string
}

func (t *Token) String() string {
	return fmt.Sprintf("type=%s value=%c", t.Type.String(), t.Char)
}

func NewToken(tokenType TokenType, char byte) *Token {
	return &Token{
		Type: tokenType,
		Char: char,
	}
}
