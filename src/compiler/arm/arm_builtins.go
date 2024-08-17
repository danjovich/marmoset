package arm

import (
	"embed"
	"fmt"
	"marmoset/code"
	"marmoset/code/arm"
	"marmoset/compiler"
)

// will be used to embed assembly files in the same directory
//
//go:embed asm
var asm embed.FS

type Builtin struct {
	Source       string
	UsedBuiltins []string // other builtins used by this one
}

var Builtins = map[string]Builtin{
	// puts a char to stdout
	"put": {
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
		UsedBuiltins: []string{},
	},
	// gets a char from stdin
	"get": {
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
		UsedBuiltins: []string{},
	},
	"putint": {
		Source:       makeAsm("putint"),
		UsedBuiltins: []string{"put"},
	},
	"putintln": {
		Source:       makeAsm("putintln"),
		UsedBuiltins: []string{"put", "putint"},
	},
}

func MakeBuiltin(index int) string {
	return Builtins[compiler.Builtins[index]].Source
}

func GetBuiltinIndex(name string) int {
	for index, builtin := range compiler.Builtins {
		if builtin == name {
			return index
		}
	}
	panic(fmt.Sprintf("unexpected error: builtin %s not found", name))
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

func makeAsm(name string) string {
	asm, err := asm.ReadFile(fmt.Sprintf("asm/%s.s", name))
	if err != nil {
		panic(fmt.Sprintf("unexpected error: %s", err))
	}

	return string(asm)
}
