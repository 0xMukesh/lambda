package evaluator

import (
	"log"

	"github.com/0xmukesh/lambda/internal/ast"
)

func isValue(node ast.Node) bool {
	_, ok := node.(*ast.Abstraction)
	return ok
}

func shift(node ast.Node, cutoff, amount int) ast.Node {
	switch n := node.(type) {
	case *ast.Identifier:
		if n.Index >= cutoff {
			return &ast.Identifier{
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
	case *ast.Identifier:
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

func Eval(node ast.Node) ast.Node {
	for {
		switch n := node.(type) {
		case *ast.Application:
			if isValue(n.Lhs) && isValue(n.Rhs) {
				abs, ok := n.Lhs.(*ast.Abstraction)
				if !ok {
					log.Fatalf("expected lhs to be abstraction, got: %s", n.Lhs.Repr())
				}

				node = betaReduction(abs.Body, n.Rhs)
			} else if isValue(n.Lhs) {
				n.Rhs = Eval(n.Rhs)
			} else {
				n.Lhs = Eval(n.Lhs)
			}
		default:
			if isValue(node) {
				return node
			}
		}
	}
}
