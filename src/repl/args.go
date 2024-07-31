package repl

import (
	"fmt"
	"os"
	"regexp"
)

const (
	VERBOSE     = "-v"
	VERBOSE_V   = "--verbose"
	INTERPRET   = "-i"
	INTERPRET_V = "--interpret"
)

type Args struct {
	Verbose        bool
	UseInterpreter bool
}

func NewArgs() *Args {
	args := &Args{}

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) > 0 {
		value := ""
		r := regexp.MustCompile(`^\-{1,2}[^\-]+`)

		for i, arg := range argsWithoutProg {
			if len(argsWithoutProg) > i+1 {
				match := r.MatchString(argsWithoutProg[i+1])
				if !match {
					value = argsWithoutProg[i+1]
				}
			}

			args.parseArg(arg, value)
		}
	}

	return args
}

func (args *Args) parseArg(arg string, value string) {
	switch arg {
	case VERBOSE, VERBOSE_V:
		assertHasNoValue(arg, value)
		args.Verbose = true
	case INTERPRET, INTERPRET_V:
		assertHasNoValue(arg, value)
		args.UseInterpreter = true
	default:
		panic(fmt.Sprintf("unknown argument %s", arg))
	}
}

func assertHasNoValue(arg string, value string) {
	if value != "" {
		panic(fmt.Sprintf("arg %s does not expect a positional argument", arg))
	}
}

// for future use
// func assertHasValue(arg string, value string) {
// 	if value == "" {
// 		panic(fmt.Sprintf("arg %s expects a value", arg))
// 	}
// }
