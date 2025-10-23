package repl

import (
	"bufio"
	"fmt"
	"gosling/evaluator"
	"gosling/lexer"
	"gosling/object"
	"gosling/parser"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

const PROMPT = "$ "

type CommandHistory struct {
	commands []string
	index    int
	maxSize  int
}

func NewCommandHistory(maxSize int) *CommandHistory {
	return &CommandHistory{
		commands: make([]string, 0),
		index:    -1,
		maxSize:  maxSize,
	}
}

func (h *CommandHistory) Add(command string) {
	if command == "" {
		return
	}

	// Don't add duplicate consecutive commands
	if len(h.commands) > 0 && h.commands[len(h.commands)-1] == command {
		h.index = -1
		return
	}

	h.commands = append(h.commands, command)

	// Keep only maxSize commands
	if len(h.commands) > h.maxSize {
		h.commands = h.commands[1:]
	}

	h.index = -1
}

func (h *CommandHistory) Previous() string {
	if len(h.commands) == 0 {
		return ""
	}

	if h.index == -1 {
		h.index = len(h.commands) - 1
	} else if h.index > 0 {
		h.index--
	}

	return h.commands[h.index]
}

func (h *CommandHistory) Next() string {
	if len(h.commands) == 0 || h.index == -1 {
		return ""
	}

	if h.index < len(h.commands)-1 {
		h.index++
		return h.commands[h.index]
	} else {
		h.index = -1
		return ""
	}
}

// New function that handles the welcome message
func StartWithWelcome(in io.Reader, out io.Writer, username string) {

	// Print welcome messages with forced line positioning
	fmt.Printf("\rWelcome to my language, Gosling\n")
	fmt.Printf("\rUser: %s\n", username)
	fmt.Printf("\rEnter Valid Gosling commands and see what happens\n")

	// Now start the REPL
	Start(in, out)
}

// Keep the original Start function for backward compatibility
func Start(in io.Reader, out io.Writer) {
	env := object.NewEnvironment()
	history := NewCommandHistory(100) // Keep last 100 commands

	// Check if we're in a terminal
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		startBasicREPL(in, out, env)
		return
	}

	// Print REPL-specific message BEFORE entering raw mode
	fmt.Printf("\rGosling REPL - Use Up/Down arrows for history, Ctrl+C to exit\n")

	// Set terminal to raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("Failed to set raw mode: %v\n", err)
		startBasicREPL(in, out, env)
		return
	}
	// defer func to check error of term.Restore
	defer func(oldState *term.State) {
		err := term.Restore(int(os.Stdin.Fd()), oldState)
		if err != nil {
			fmt.Printf("Failed to restore terminal state: %v\n", err)
		}
	}(oldState)

	for {
		// Force cursor to column 0 and print prompt
		fmt.Printf("\r\033[2K%s", PROMPT)
		line := readLineWithHistory(history)

		if line == "" {
			continue
		}

		// Handle exit commands
		if line == "\\exit" || line == "\\quit" {
			fmt.Printf("\nGoodbye!\n")
			break
		}

		history.Add(line)

		// Parse and evaluate - force new line before output
		fmt.Printf("\n")

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)

		if evaluated != nil {
			fmt.Printf("%s\n", evaluated.Inspect())
		}
	}
}

func readLineWithHistory(history *CommandHistory) string {
	var line strings.Builder

	for {
		var buf [1]byte
		n, err := os.Stdin.Read(buf[:])
		if err != nil || n == 0 {
			continue
		}

		b := buf[0]

		switch b {
		case 3: // Ctrl+C
			fmt.Printf("\nGoodbye!\n")
			os.Exit(0)
		case 13: // Enter (CR)
			return line.String()
		case 10: // Line Feed (LF)
			return line.String()
		case 127: // Backspace (DEL on Mac)
			if line.Len() > 0 {
				current := line.String()
				line.Reset()
				if len(current) > 0 {
					line.WriteString(current[:len(current)-1])
				}
				fmt.Printf("\b \b")
			}
		case 8: // Backspace (BS)
			if line.Len() > 0 {
				current := line.String()
				line.Reset()
				if len(current) > 0 {
					line.WriteString(current[:len(current)-1])
				}
				fmt.Printf("\b \b")
			}
		case 27: // Escape sequence (arrow keys)
			// Read the next two bytes for arrow keys
			var seq [2]byte
			n, _ := os.Stdin.Read(seq[:])
			if n >= 2 && seq[0] == 91 { // '['
				switch seq[1] {
				case 65: // Up arrow
					prev := history.Previous()
					clearLineAndRedraw(prev)
					line.Reset()
					line.WriteString(prev)
				case 66: // Down arrow
					next := history.Next()
					clearLineAndRedraw(next)
					line.Reset()
					line.WriteString(next)
				}
			}
		default:
			if b >= 32 && b <= 126 { // Printable ASCII characters
				line.WriteByte(b)
				fmt.Printf("%c", b)
			}
		}
	}
}

func clearLineAndRedraw(newText string) {
	// Move to beginning of line and clear it completely
	fmt.Printf("\r\033[2K%s%s", PROMPT, newText)
}

func startBasicREPL(in io.Reader, out io.Writer, env *object.Environment) {
	_, err := fmt.Fprintf(out, "Gosling REPL (basic mode)\n")
	if err != nil {
		fmt.Printf("Failed to start REPL\n")
		return
	}
	scanner := bufio.NewScanner(in)

	for {
		_, err := fmt.Fprint(out, PROMPT)
		if err != nil {
			fmt.Printf("Failed to print prompt in REPL\n")
			return
		}
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)

		if evaluated != nil {
			_, err := fmt.Fprintf(out, "%s\n", evaluated.Inspect())
			if err != nil {
				fmt.Printf("Failed to print prompt in REPL\n")
				return
			}

		}
	}
}

func printParseErrors(errors []string) {
	for _, msg := range errors {
		fmt.Printf("\t%s\n", msg)
	}
}
