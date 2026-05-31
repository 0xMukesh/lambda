package ast

import (
	"fmt"
)

type Node interface {
	Repr() string
}
type Ast []Node

type Variable struct {
	Value byte
	Index int
}

func (id *Variable) Repr() string {
	return fmt.Sprintf("%c", id.Value)
}

type Abstraction struct {
	Param string
	Body  Node
}

func (ab *Abstraction) Repr() string {
	return fmt.Sprintf("\\%s.%s", ab.Param, ab.Body.Repr())
}

type Application struct {
	Lhs Node
	Rhs Node
}

func (ap *Application) Repr() string {
	return fmt.Sprintf("(%s %s)", ap.Lhs.Repr(), ap.Rhs.Repr())
}

type Assignment struct {
	Name string
	Body Node
}

func (as *Assignment) Repr() string {
	return fmt.Sprintf("%s = %s", as.Name, as.Body.Repr())
}

type NamedTermRef struct {
	Name string
}

func (nt *NamedTermRef) Repr() string {
	return nt.Name
}
