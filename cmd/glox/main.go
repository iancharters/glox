package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/iancharters/glox/parser"
	"github.com/iancharters/glox/scanner"
)

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	}

	if len(args) == 1 {
		runFile(args[0])
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err := run(string(data)); err != nil {
		panic(err)
	}

}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		if input == "\n" {
			break
		}

		if err := run(input); err != nil {
			fmt.Println(err)
		}
	}
}

func run(source string) error {
	s := scanner.New(source)

	tokens, err := s.ScanTokens()
	if err != nil {
		panic(err)
	}

	p := parser.New(tokens)

	expression, err := p.Parse()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	astprinter := parser.Printer{}
	printed := astprinter.Print(expression)

	fmt.Println(printed)

	return nil
}
