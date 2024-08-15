package args

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
	Program        string
}

func NewArgs() (*Args, error) {
	args := &Args{}

	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) > 0 {
		value := ""
		r := regexp.MustCompile(`^\-{1,2}[^\-]+`)

		lastIndex := len(argsWithoutProg) - 1
		if !r.MatchString(argsWithoutProg[lastIndex]) {
			args.Program = argsWithoutProg[lastIndex]
			argsWithoutProg = argsWithoutProg[:lastIndex]
		}

		for i, arg := range argsWithoutProg {
			if len(argsWithoutProg) > i+1 {
				match := r.MatchString(argsWithoutProg[i+1])
				if !match {
					value = argsWithoutProg[i+1]
				}
			}

			err := args.parseArg(arg, value)
			if err != nil {
				return nil, err
			}
		}
	}

	return args, nil
}

func (args *Args) parseArg(arg string, value string) error {
	switch arg {
	case VERBOSE, VERBOSE_V:
		err := assertHasNoValue(arg, value)
		if err != nil {
			return err
		}
		args.Verbose = true
	case INTERPRET, INTERPRET_V:
		err := assertHasNoValue(arg, value)
		if err != nil {
			return err
		}
		args.UseInterpreter = true
	default:
		return fmt.Errorf("unknown argument %s", arg)
	}
	return nil
}

func assertHasNoValue(arg string, value string) error {
	if value != "" {
		return fmt.Errorf("arg %s does not expect a positional argument", arg)
	}
	return nil
}

// for future use
// func assertHasValue(arg string, value string) {
// 	if value == "" {
// 		panic(fmt.Sprintf("arg %s expects a value", arg))
// 	}
// }
