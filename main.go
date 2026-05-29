package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Expr interface {
	ExprRepr() string
}

// variables - x, y, z
type Variable struct {
	Name string
}

func (v *Variable) ExprRepr() string {
	return v.Name
}

// functions - \x.x
type Function struct {
	Arg  Variable
	Body Expr
}

func (f *Function) ExprRepr() string {
	return fmt.Sprintf("\\%s.%s", f.Arg.ExprRepr(), f.Body.ExprRepr())
}

// function application - (\x.x)(a) => a
type Application struct {
	Lhs Expr
	Rhs Expr
}

func (a *Application) ExprRepr() string {
	return fmt.Sprintf("%s %s", a.Lhs.ExprRepr(), a.Rhs.ExprRepr())
}

func main() {
	for {
		fmt.Printf("λ [in]: ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("failed to read input: %s", err)
		}

		if line == "q\n" || line == "exit\n" {
			break
		}

		fmt.Printf("λ [out]: %s", line)
	}
}
