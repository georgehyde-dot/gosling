package repl

import (
	"bufio"
	"fmt"
	"gosling/lexer"
	"gosling/parser"
	"gosling/token"
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

		n, err := io.WriteString(out, program.String())
		if err != nil {
			fmt.Printf("error writing to out, %d bytes\n", n)
		}
		n, err = io.WriteString(out, "\n")
		if err != nil {
			fmt.Printf("error writing to out, %d bytes\n", n)
		}

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
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
