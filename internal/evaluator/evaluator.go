package evaluator

import (
	"github.com/0xmukesh/lambda/internal/ast"
)

func shift(node ast.Node, cutoff, amount int) ast.Node {
	switch n := node.(type) {
	case *ast.Identifier:
		if n.Index >= cutoff {
			return &ast.Identifier{
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
		case *ast.Abstraction:
			return &ast.Abstraction{
				Param: n.Param,
				Body:  Eval(n.Body),
			}
		case *ast.Application:
			n.Lhs = Eval(n.Lhs)

			abs, ok := n.Lhs.(*ast.Abstraction)
			if !ok {
				return &ast.Application{
					Lhs: n.Lhs,
					Rhs: Eval(n.Rhs),
				}
			}

			node = betaReduction(abs.Body, n.Rhs)
		default:
			return node
		}
	}
}
