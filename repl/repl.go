package repl

import (
	"bufio"
	"fmt"
	"gosling/evaluator"
	"gosling/lexer"
	"gosling/parser"
	"io"
)

const PROMPT = "$"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.LexRepl(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program)

		if evaluated != nil {
			_, err := io.WriteString(out, evaluated.Inspect())
			if err != nil {
				fmt.Print("Failed to write evaluated string")
			}
			_, err = io.WriteString(out, "\n")
			if err != nil {
				fmt.Print("Failed to write new line")
			}
		}
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		n, err := io.WriteString(out, "\t"+msg+"\n")
		if err != nil {
			fmt.Printf("error writing to out, %d bytes\n", n)
		}
	}
}
