package main

import (
	"fmt"
	"io"

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
		toks, err := lexer.LexAll()
		if err != nil {
			panic(fmt.Sprintf("failed to lex: %s", err))
		}

		parser := parser.NewParser(toks)
		ast, err := parser.Parse()
		if err != nil {
			panic(fmt.Sprintf("failed to parse: %s", err))
		}

		result := evaluator.Eval(ast)

		fmt.Printf("%s\n", result.Repr())
	}
}
