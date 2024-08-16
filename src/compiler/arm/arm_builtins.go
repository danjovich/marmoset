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
	{Name: "put",
		// makes the syscall, restores lr and returns
		Source: fmt.Sprintf(`put:
%s
L0_put: @put
%s

%s
`,
			arm.MakeFunctionPreamble(1*4),
			MakeSyscall(4, `mov r0, #1
	add r1, fp, #-4
	mov r2, #1`), makeReturn("put", 2)),
	},
}

func MakeBuiltin(index int) string {
	return Builtins[index].Source
}

func makeReturn(name string, lrIndex int) string {
	returnCode, err := arm.Make(code.OpReturn, 3, name, lrIndex)
	if err != nil {
		panic(fmt.Sprintf("unexpected error: %s", err))
	}

	return returnCode
}
