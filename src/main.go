package main

import (
	"fmt"
	"marmoset/args"
	"marmoset/repl"
	"os"
	"os/user"
)

func main() {
	args := args.New()

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the Marmoset programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands!\n")
	repl.Start(os.Stdin, os.Stdout, args)
}
