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
	args, err := args.NewArgs()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	if args.Program != "" {
		source, err := os.ReadFile(args.Program)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error when reading file: %s", err)
			os.Exit(1)
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
			printParserErrors(os.Stderr, p.Errors())
			os.Exit(1)
		}

		comp := compiler.NewWithState(symbolTable, constants)
		err = comp.Compile(program)
		if err != nil {
			fmt.Fprintf(os.Stderr, "compilation failed:\n %s\n", err)
			os.Exit(1)
		}

		arm_compiler := arm.New(comp)
		err = arm_compiler.Compile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "compilation failed:\n %s\n", err)
			os.Exit(1)
		}

		return
	}

	user, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "user fetching failed:\n %s\n", err)
		os.Exit(1)
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
