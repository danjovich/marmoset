package main

import (
	"fmt"
	"marmoset/repl"
	"os"
	"os/user"
)

func main() {
	args := repl.NewArgs()

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the Marmoset programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout, args)
}
