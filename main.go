package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/0xmukesh/lambda/internal/ast"
	"github.com/0xmukesh/lambda/internal/evaluator"
	"github.com/0xmukesh/lambda/internal/lexer"
	"github.com/0xmukesh/lambda/internal/parser"
	"github.com/chzyer/readline"
)

func main() {
	if len(os.Args) < 2 {
		panic("invalid usage. correct usage: lambda <command>. available commands: repl, run")
	}

	cmd := os.Args[1]
	defs := make(map[string]ast.Node)

	if cmd == "repl" {
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

			runPipeline([]byte(line), defs)
		}
	} else if cmd == "run" {
		var file string
		runCmd := flag.NewFlagSet("run", flag.ExitOnError)
		runCmd.StringVar(&file, "file", "", "path to the file which contains the source code")
		runCmd.Parse(os.Args[2:])

		if file == "" {
			panic("missing --file flag")
		}

		ext := filepath.Ext(file)
		if ext != ".lm" {
			panic(fmt.Errorf("file extension must be `.lm`, got %s", ext))
		}

		source, err := os.ReadFile(file)
		if err != nil {
			panic(fmt.Errorf("failed to read %s file: %s", file, err))
		}

		content := strings.TrimSpace(string(source))
		for line := range strings.SplitSeq(content, "\n") {
			runPipeline([]byte(line), defs)
		}
	} else {
		panic(fmt.Sprintf("unknown command: %s", cmd))
	}
}

func runPipeline(source []byte, defs map[string]ast.Node) {
	lexer := lexer.NewLexer(source)
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
		result := evaluator.Eval(expanded, defs)
		if result != nil {
			fmt.Println(result.Repr())
		}
	}
}
