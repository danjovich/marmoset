package main

import (
	"fmt"
	"io"
	"marmoset/args"
	"marmoset/compiler"
	"marmoset/compiler/arm"
	"marmoset/lexer"
	"marmoset/object"
	"marmoset/parser"
	"marmoset/repl"
	"os"
	"os/user"
)

func main() {
	args := args.NewArgs()

	if args.Program != "" {
		source, err := os.ReadFile(args.Program)
		if err != nil {
			fmt.Printf("error when reading file: %s", err)
			return
		}

		constants := []object.Object{}
		symbolTable := compiler.NewSymbolTable()
		for i, v := range object.Builtins {
			symbolTable.DefineBuiltin(i, v.Name)
		}

		l := lexer.New(string(source))
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(os.Stdout, p.Errors())
			return
		}

		comp := compiler.NewWithState(symbolTable, constants)
		err = comp.Compile(program)
		if err != nil {
			fmt.Printf("compilation failed:\n %s\n", err)
			return
		}

		arm_compiler := arm.New(comp)
		err = arm_compiler.Compile()
		if err != nil {
			fmt.Printf("compilation failed:\n %s\n", err)
			return
		}

		return
	}

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the Marmoset programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout, args)
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
