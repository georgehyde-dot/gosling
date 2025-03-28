package repl

import (
	"bufio"
	"fmt"
	"gosling/lexer"
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
		fmt.Fprintln(out, line)
		l := lexer.LexRepl(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
