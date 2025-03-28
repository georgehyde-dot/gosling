package main

import (
	"fmt"
	"gosling/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Welcome to my language, Gosling\nUser: %s\n", user.Username)
	fmt.Printf("Enter Valid Gosling commands and see what happens\n")
	repl.Start(os.Stdin, os.Stdout)
}
