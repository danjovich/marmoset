package main

import (
	"fmt"
	"io"
	"marmoset/args"
	"marmoset/compiler"
	"marmoset/compiler/arm"
	"marmoset/lexer"
	"marmoset/parser"
	"os"
)

func main() {
	args, err := args.NewArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	source, err := os.ReadFile(args.Program)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error when reading file: %s\n", err)
		os.Exit(1)
	}

	l := lexer.New(string(source))
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(os.Stderr, p.Errors())
		os.Exit(1)
	}

	comp := compiler.New()
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
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
