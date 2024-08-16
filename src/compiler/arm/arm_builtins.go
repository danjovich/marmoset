package arm

import (
	"fmt"
	"marmoset/code"
	"marmoset/code/arm"
)

type Builtin struct {
	Name   string
	Source string
}

var Builtins = []Builtin{
	// puts a char to stdout
	{Name: "put",
		// makes the syscall, restores lr and returns
		Source: fmt.Sprintf(`put:
%s
L0_put: @put
%s

%s
`,
			// one arg
			arm.MakeFunctionPreamble(1*4),
			// args for write (syscall 4) are stdout (1), location of the char (fp - 1 == first argument)
			// and size 1
			MakeSyscall(4, `mov r0, #1
	add r1, fp, #-4
	mov r2, #1`), makeReturn("put", 2, false)),
	},
	// gets a char from stdin
	{Name: "get",
		// makes the syscall, restores lr and returns
		Source: fmt.Sprintf(`get:
%s
L0_get: @get
%s

%s
`,
			// no args
			arm.MakeFunctionPreamble(0),
			// args for read (syscall 3) are stdin (0), location of the char (sp + 1 == next sp value)
			// and size 2 (to have space for the \n)
			MakeSyscall(3, `mov r0, #0
	sub sp, sp, #4
	mov r1, sp
	mov r2, #2`), makeReturn("get", 2, true)),
	},
}

func MakeBuiltin(index int) string {
	return Builtins[index].Source
}

func makeReturn(name string, lrIndex int, isReturnValue bool) string {
	var op code.Opcode

	if isReturnValue {
		op = code.OpReturnValue
	} else {
		op = code.OpReturn
	}

	returnCode, err := arm.Make(op, 3, name, lrIndex)
	if err != nil {
		panic(fmt.Sprintf("unexpected error: %s", err))
	}

	return returnCode
}
