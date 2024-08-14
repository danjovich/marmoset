package arm

import (
	"fmt"
	"marmoset/code"
	"marmoset/code/arm"
	"marmoset/compiler"
	"marmoset/object"
	"strconv"
)

func Compile(compiler *compiler.Compiler) error {
	// TODO: use arm-linux-gnueabihf-as?

	bytecode := compiler.Bytecode()
	symbols := compiler.SymbolTable

	return compileFromInstructionsAndSymbols(bytecode, symbols)
}

func compileFromInstructionsAndSymbols(bytecode *compiler.Bytecode, symbols *compiler.SymbolTable) error {
	ins := bytecode.Instructions

	for ip := 0; ip < len(ins); ip++ {
		op := code.Opcode(ins[ip])
		index := ip
		args := []interface{}{}

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(ins[ip+1:])
			constant := bytecode.Constants[constIndex]
			arg, err := generateConstantArg(constant)
			if err != nil {
				return err
			}
			args = append(args, arg)
			ip += 2

		case code.OpJump, code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(ins[ip+1:]))
			ip += 2
			arg := fmt.Sprintf("L%d", pos)
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
			args = append(args, numArgs)

		case code.OpSetLocal, code.OpGetLocal:
			localIndex := code.ReadUint16(ins[ip+1:])
			ip += 2
			localName, ok := symbols.ResolveName(int(localIndex), compiler.LocalScope)
			if !ok {
				return fmt.Errorf("local of index %d not found", localIndex)
			}
			args = append(args, localName)

		case code.OpGetBuiltin:
			builtinIndex := code.ReadUint16(ins[ip+1:])
			ip += 2
			builtinName, ok := symbols.ResolveName(int(builtinIndex), compiler.BuiltinScope)
			if !ok {
				return fmt.Errorf("builtin of index %d not found", builtinIndex)
			}
			args = append(args, builtinName)
		}

		asm, err := arm.Make(op, index, args)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", asm)
	}

	return nil
}

// TODO: test this function
func generateConstantArg(constant object.Object) ([]string, error) {
	switch constant.(type) {
	case *object.CompiledFunction:
		compiledFunction := constant.(*object.CompiledFunction)
		result := fmt.Sprintf("=%s", compiledFunction.Name)
		return []string{result}, nil

	case *object.Integer:
		integer := constant.(*object.Integer)
		result := fmt.Sprintf("#%d", integer.Value)
		return []string{result}, nil

	case *object.String:
		value := constant.(*object.String).Value
		asciiValues, err := strconv.Atoi(value)
		if err != nil {
			return []string{}, err
		}
		result := []byte{}
		for asciiValues != 0 {
			asciiValue := asciiValues % 0xFF
			result = fmt.Appendf(result, "#%d", asciiValue)
			asciiValues <<= 8
		}
	}

	return []string{}, fmt.Errorf("invalid type %T for constant", constant.Type())
}
