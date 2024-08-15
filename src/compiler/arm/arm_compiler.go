package arm

import (
	"fmt"
	"marmoset/code"
	arm "marmoset/code/arm"
	"marmoset/compiler"
	"marmoset/object"
	"slices"
	"strconv"
)

type ArmCompiler struct {
	compiler *compiler.Compiler
}

func New(compiler *compiler.Compiler) *ArmCompiler {
	return &ArmCompiler{
		compiler: compiler,
	}
}

func (ac *ArmCompiler) Compile() error {
	// TODO: use arm-linux-gnueabihf-as?

	constants := ac.compiler.Constants
	globalFunctions := []string{}

	for _, scope := range ac.compiler.AllScopes {
		fmt.Printf("%s:\n", scope.Name)
		if scope.IsMain {
			fmt.Print("	mov sp, #0x4000\n	mov fp, sp\n\n")
		}

		globalFunctions = append(globalFunctions, scope.Name)

		err := compileFromInstructionsAndSymbols(*scope, constants)
		if err != nil {
			return err
		}

		if scope.IsMain {
			fmt.Print("_end: b _end\n\n")
		}
	}

	for _, name := range ac.compiler.SymbolTable.GetAllGlobalNames() {
		if !slices.Contains(globalFunctions, name) {
			fmt.Printf("_%s: .word 0x0\n", name)
		} else {
			fmt.Printf("_%s: .word %s\n", name, name)
		}
	}

	return nil
}

func compileFromInstructionsAndSymbols(scope compiler.CompilationScope, constants []object.Object) error {
	ins := scope.Instructions
	symbols := scope.SymbolTable

	startOfTheStack := 0
	scopeName := fmt.Sprintf("_%s", scope.Name)
	if scope.IsMain {
		scopeName = ""
	} else {
		// preamble -> must open space for ALL local variables, including arguments and old lr
		startOfTheStack = (len(symbols.GetAllLocalNames()) + 1)
		fmt.Printf("	sub sp, sp, #%d\n\n", startOfTheStack*4)
	}

	for ip := 0; ip < len(ins); ip++ {
		op := code.Opcode(ins[ip])
		index := ip
		args := []interface{}{}

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(ins[ip+1:])
			constant := constants[constIndex]
			arg, err := generateConstantArgs(constant)
			if err != nil {
				return err
			}
			args = append(args, arg...)
			ip += 2

		case code.OpJump, code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(ins[ip+1:]))
			ip += 2
			arg := fmt.Sprintf("L%d%s", pos, scopeName)
			args = append(args, arg)

		case code.OpGetGlobal, code.OpSetGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			ip += 2
			globalName, ok := symbols.ResolveName(int(globalIndex), compiler.GlobalScope)
			if !ok {
				return fmt.Errorf("global of index %d not found", globalIndex)
			}
			args = append(args, globalName)

		case code.OpArray:
			numElements := int(code.ReadUint16(ins[ip+1:]))
			ip += 2
			args = append(args, numElements)

		case code.OpCall:
			numArgs := code.ReadUint8(ins[ip+1:])
			args = append(args, int(numArgs))
			ip += 1

		case code.OpSetLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			ip += 1
			args = append(args, int(localIndex))

		case code.OpGetLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			ip += 1
			args = append(args, int(localIndex))

		case code.OpGetBuiltin:
			builtinIndex := code.ReadUint8(ins[ip+1:])
			ip += 1
			builtinName, ok := symbols.ResolveName(int(builtinIndex), compiler.BuiltinScope)
			if !ok {
				return fmt.Errorf("builtin of index %d not found", builtinIndex)
			}
			args = append(args, builtinName)

		case code.OpReturn, code.OpReturnValue:
			args = append(args, startOfTheStack)
		}

		asm, err := arm.Make(op, index, scopeName, args...)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", asm)
	}

	return nil
}

// TODO: test this function
func generateConstantArgs(constant object.Object) ([]interface{}, error) {
	switch constant := constant.(type) {
	case *object.CompiledFunction:
		result := fmt.Sprintf("#%s", constant.Name)
		return []interface{}{result}, nil

	case *object.Integer:
		result := fmt.Sprintf("#%d", constant.Value)
		return []interface{}{result}, nil

	case *object.String:
		asciiValues, err := strconv.Atoi(constant.Value)
		if err != nil {
			return []interface{}{}, err
		}
		result := []interface{}{}
		for asciiValues != 0 {
			asciiValue := asciiValues % 0xFF
			result = append(result, fmt.Sprintf("#%d", asciiValue))
			asciiValues <<= 8
		}
		return result, nil
	}

	return []interface{}{}, fmt.Errorf("invalid type %T for constant", constant.Type())
}
