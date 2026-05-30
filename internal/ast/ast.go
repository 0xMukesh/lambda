package ast

import (
	"fmt"
)

type Node interface {
	Repr() string
}
type Ast []Node

// identifiers - x, y, z
type Identifier struct {
	Value string
}

func (id *Identifier) Repr() string {
	return id.Value
}

// abstraction - \x.x
type Abstraction struct {
	Param *Identifier
	Body  Node
}

func (ab *Abstraction) Repr() string {
	return fmt.Sprintf("(abs %s %s)", ab.Param.Repr(), ab.Body.Repr())
}

// application - (\x.x) a => a
type Application struct {
	Lhs Node
	Rhs Node
}

func (ap *Application) Repr() string {
	return fmt.Sprintf("(app %s %s)", ap.Lhs.Repr(), ap.Rhs.Repr())
}
