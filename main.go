package main

import (
	"fmt"
	"io"

	"github.com/0xmukesh/lambda/internal/ast"
	"github.com/0xmukesh/lambda/internal/evaluator"
	"github.com/0xmukesh/lambda/internal/lexer"
	"github.com/0xmukesh/lambda/internal/parser"
	"github.com/chzyer/readline"
)

func main() {
	rl, err := readline.New("λ: ")
	if err != nil {
		panic(fmt.Sprintf("failed to setup readline: %s", err))
	}
	defer rl.Close()

	defs := make(map[string]ast.Node)

	for {
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt || err == io.EOF {
				fmt.Println("exiting...")
				break
			} else {
				panic(fmt.Sprintf("failed to read line: %s", err))
			}
		}

		if line == "" {
			continue
		}

		lexer := lexer.NewLexer([]byte(line))
		tokens, err := lexer.LexAll()
		if err != nil {
			panic(fmt.Sprintf("failed to lex: %s", err))
		}

		parser := parser.NewParser(tokens)
		node, err := parser.Parse()
		if err != nil {
			panic(fmt.Sprintf("failed to parse: %s", err))
		}

		switch n := node.(type) {
		case *ast.Assignment:
			defs[n.Name] = n.Body
		default:
			expanded := evaluator.Expand(node, defs)
			result := evaluator.Eval(expanded)
			fmt.Println(result.Repr())
		}
	}
}
