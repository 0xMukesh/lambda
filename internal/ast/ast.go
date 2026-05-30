package ast

import (
	"fmt"
)

type Node interface {
	Repr() string
}
type Ast []Node

type Identifier struct {
	Value byte
	Index int
}

func (id *Identifier) Repr() string {
	return fmt.Sprintf("%c", id.Value)
}

type Abstraction struct {
	Param string
	Body  Node
}

func (ab *Abstraction) Repr() string {
	return fmt.Sprintf("(λ%s.%s)", ab.Param, ab.Body.Repr())
}

type Application struct {
	Lhs Node
	Rhs Node
}

func (ap *Application) Repr() string {
	return fmt.Sprintf("(%s %s)", ap.Lhs.Repr(), ap.Rhs.Repr())
}
