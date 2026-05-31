package evaluator

import (
	"fmt"

	"github.com/0xmukesh/lambda/internal/ast"
)

func shift(node ast.Node, cutoff, amount int) ast.Node {
	switch n := node.(type) {
	case *ast.Variable:
		if n.Index >= cutoff {
			return &ast.Variable{
				Value: n.Value,
				Index: n.Index + amount,
			}
		}

		return n
	case *ast.Abstraction:
		return &ast.Abstraction{
			Param: n.Param,
			Body:  shift(n.Body, cutoff+1, amount),
		}
	case *ast.Application:
		return &ast.Application{
			Lhs: shift(n.Lhs, cutoff, amount),
			Rhs: shift(n.Rhs, cutoff, amount),
		}
	}

	return node
}

func substitute(body ast.Node, target int, arg ast.Node) ast.Node {
	switch n := body.(type) {
	case *ast.Variable:
		if n.Index == target {
			return arg
		}

		return n
	case *ast.Abstraction:
		return &ast.Abstraction{
			Param: n.Param,
			Body:  substitute(n.Body, target+1, shift(arg, 0, 1)),
		}
	case *ast.Application:
		return &ast.Application{
			Lhs: substitute(n.Lhs, target, arg),
			Rhs: substitute(n.Rhs, target, arg),
		}
	}

	return body
}

func betaReduction(body, arg ast.Node) ast.Node {
	shifted := shift(arg, 0, 1)
	substituted := substitute(body, 0, shifted)
	return shift(substituted, 0, -1)
}

func Expand(node ast.Node, defs map[string]ast.Node) ast.Node {
	switch n := node.(type) {
	case *ast.NamedTermRef:
		def, ok := defs[n.Name]
		if !ok {
			panic(fmt.Sprintf("unknown named term: %s", n.Name))
		}

		return Expand(def, defs)
	case *ast.Abstraction:
		return &ast.Abstraction{
			Param: n.Param,
			Body:  Expand(n.Body, defs),
		}
	case *ast.Application:
		return &ast.Application{
			Lhs: Expand(n.Lhs, defs),
			Rhs: Expand(n.Rhs, defs),
		}
	}

	return node
}

func Eval(node ast.Node, defs map[string]ast.Node) ast.Node {
	for {
		switch n := node.(type) {
		case *ast.NamedTermRef:
			node = defs[n.Name]
		case *ast.Abstraction:
			return &ast.Abstraction{
				Param: n.Param,
				Body:  Eval(n.Body, defs),
			}
		case *ast.Application:
			n.Lhs = Eval(n.Lhs, defs)
			abs, ok := n.Lhs.(*ast.Abstraction)
			if !ok {
				return &ast.Application{
					Lhs: n.Lhs,
					Rhs: Eval(n.Rhs, defs),
				}
			}

			node = betaReduction(abs.Body, n.Rhs)
		default:
			return node
		}
	}
}
