package object

import (
	"fmt"
)

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	// print char to STDOUT
	{"put",
		&Builtin{
			Fn: func(args ...Object) Object {
				fmt.Print(args[0].Inspect())
				return nil
			},
		},
	},
	{"get",
		&Builtin{
			Fn: func(args ...Object) Object {
				// reader := bufio.NewReader(os.Stdin)
				// input, _ := reader.ReadString('\n')
				// return &Integer{Value: input}
				return nil
			},
		},
	},
	{"putint",
		&Builtin{
			Fn: func(args ...Object) Object {
				fmt.Print(args[0])
				return nil
			},
		},
	},
	{"putintln",
		&Builtin{
			Fn: func(args ...Object) Object {
				fmt.Println(args[0])
				return nil
			},
		},
	},
}
