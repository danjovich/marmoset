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

	symbols := ac.compiler.SymbolTable
	constants := ac.compiler.Constants
	globalFunctions := []string{}

	for _, name := range symbols.GetAllGlobalNames() {
		if !slices.Contains(globalFunctions, name) {
			fmt.Printf("_%s: .word 0x0\n", name)
		} else {
			fmt.Printf("_%s: .word %s\n", name, name)
		}
	}

	fmt.Printf("\n")

	for _, scope := range ac.compiler.AllScopes {
		fmt.Printf("%s:\n", scope.Name)
		globalFunctions = append(globalFunctions, scope.Name)

		err := compileFromInstructionsAndSymbols(*scope, constants, symbols)
		if err != nil {
			return err
		}
	}

	return nil
}

func compileFromInstructionsAndSymbols(scope compiler.CompilationScope, constants []object.Object, symbols *compiler.SymbolTable) error {
	ins := scope.Instructions
	// locals store declared local variables and counts the amount of changes in the stack since
	// their declaration to allow for knowing where they are when they are read
	locals := make(map[string]int)

	scopeName := fmt.Sprintf("_%s", scope.Name)
	if scope.IsMain {
		scopeName = ""
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
			localName, ok := symbols.ResolveName(int(localIndex), compiler.LocalScope)
			if !ok {
				return fmt.Errorf("local of index %d not found", localIndex)
			}
			// (re)starts the counting of changes in the stack since declaration
			locals[localName] = 0

		case code.OpGetLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			ip += 1
			localName, ok := symbols.ResolveName(int(localIndex), compiler.LocalScope)
			if !ok {
				return fmt.Errorf("local of index %d not found", localIndex)
			}
			stackChanges := locals[localName]
			// the amount of stack changes is passed for the generated code to be able to calculate
			// how far from the current stack pointer the local variable is
			args = append(args, stackChanges)

		case code.OpGetBuiltin:
			builtinIndex := code.ReadUint8(ins[ip+1:])
			ip += 1
			builtinName, ok := symbols.ResolveName(int(builtinIndex), compiler.BuiltinScope)
			if !ok {
				return fmt.Errorf("builtin of index %d not found", builtinIndex)
			}
			args = append(args, builtinName)
		}

		asm, changesInStack, err := arm.Make(op, index, scopeName, args...)
		if err != nil {
			return err
		}
		// adds the amount of stack changes to every previously declared local symbol count
		sumToAllValuesInMap(&locals, changesInStack)

		fmt.Printf("%s\n", asm)
	}

	return nil
}

// TODO: test this function
func generateConstantArgs(constant object.Object) ([]interface{}, error) {
	switch constant := constant.(type) {
	case *object.CompiledFunction:
		result := fmt.Sprintf("#_%s", constant.Name)
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

func sumToAllValuesInMap(stringIntMap *map[string]int, value int) {
	for key := range *stringIntMap {
		(*stringIntMap)[key] += value
	}
}
