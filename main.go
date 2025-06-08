package main

import (
	"gosling/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	// Pass the user info to the REPL to handle printing
	repl.StartWithWelcome(os.Stdin, os.Stdout, user.Username)
}
