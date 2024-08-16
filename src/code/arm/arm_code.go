package arm

import (
	"fmt"
	"marmoset/code"
)

// returns the AArch32 assembly instructions corresponding to the opcode, ip, scope and operands given
func Make(op code.Opcode, index int, scopeName string, operands ...any) (string, error) {
	label := fmt.Sprintf("L%d%s", index, scopeName)
	switch op {
	case code.OpConstant:
		if len(operands) < 1 {
			return "", fmt.Errorf("OpConstant should have at least one operand")
		}
		result := []byte(fmt.Sprintf(`%s:  @OpConstant
`, label))
		for _, operand := range operands {
			result = fmt.Appendf(result, `	%s
	push {r0}
`, operand)
		}
		return string(result), nil

	case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
		if len(operands) != 0 {
			return "", fmt.Errorf("binary operations should not have any operands")
		}
		return makeBinaryOperation(op, label), nil

	case code.OpPop:
		// TODO: pop of arrays (maybe an "array pointer", similar to stack pointer, or check previous instruction on compiler)
		if len(operands) != 0 {
			return "", fmt.Errorf("OpPop should not have any operands")
		}
		return fmt.Sprintf(`%s:  @OpPop
	pop {r0}
`, label), nil

	case code.OpTrue, code.OpFalse:
		if len(operands) != 0 {
			return "", fmt.Errorf("booleans should not have any operands")
		}
		return makeBoolean(op == code.OpTrue, label), nil

	case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
		if len(operands) != 0 {
			return "", fmt.Errorf("comparisons should not have any operands")
		}
		return makeComparison(op, label), nil

	case code.OpMinus:
		if len(operands) != 0 {
			return "", fmt.Errorf("OpMinus should not have any operands")
		}
		return fmt.Sprintf(`%s:  @OpMinus
	pop {r0}
	sub r0, #0, r0
	push {r0}
`, label), nil

	case code.OpBang:
		if len(operands) != 0 {
			return "", fmt.Errorf("OpBang should not have any operands")
		}
		return fmt.Sprintf(`%s:  @OpBang
	mov r1, #0
	pop {r0}
	cmp r0, #0
	moveq r1, #1
	push {r1}
`, label), nil

	case code.OpJumpNotTruthy:
		if len(operands) != 1 {
			return "", fmt.Errorf("OpJumpNotTruthy should have only one operand")
		}
		dest := operands[0]
		return fmt.Sprintf(`%s:  @OpJumpNotTruthy
	pop {r0}
	cmp r0, #0
	beq %s
`, label, dest), nil

	case code.OpJump:
		if len(operands) != 1 {
			return "", fmt.Errorf("OpJump should have only one operand")
		}
		dest := operands[0]
		return fmt.Sprintf(`%s:  @OpJump
	b %s
`, label, dest), nil

	case code.OpNull:
		if len(operands) != 0 {
			return "", fmt.Errorf("OpNull should not have any operands")
		}
		return fmt.Sprintf(`%s:  @OpNull
	mov r0, #0
	push {r0}
`, label), nil

	case code.OpGetGlobal:
		if len(operands) != 1 {
			return "", fmt.Errorf("OpGetGlobal should have only one operand")
		}
		name := operands[0]
		return fmt.Sprintf(`%s:  @OpGetGlobal
	ldr r0, =_%s
	ldr r1, [r0]
	push {r1}
`, label, name), nil

	case code.OpSetGlobal:
		if len(operands) != 1 {
			return "", fmt.Errorf("OpSetGlobal should have only one operand")
		}
		name := operands[0]
		return fmt.Sprintf(`%s:  @OpSetGlobal
	pop {r0}
	ldr r1, =_%s
	str r0, [r1]
`, label, name), nil

	case code.OpArray:
		if len(operands) != 1 {
			return "", fmt.Errorf("OpArray should have only one operand")
		}
		length := operands[0]
		// arrays must push their length and memory location
		return fmt.Sprintf(`%s:  @OpArray
	mov r0, #%d
	push {r0}
	sub r0, sp, #4
	push {r0}
`, label, length), nil

	case code.OpIndex:
		if len(operands) != 0 {
			return "", fmt.Errorf("OpIndex should not have any operands")
		}
		// r3 = location of array
		// r0 = index, r1 = length
		// r2 = length - index - 1 // because its zero-indexed
		// r0 = mem[sp + r2 * 4] // same as r0 = arr[index]
		// return r0
		return fmt.Sprintf(`%s:  @OpIndex
	pop {r3} 
	ldr r0, [r3, #4]
	ldr r1, [r3, #8]
	sub r2, r1, r0
	sub r2, r2, #1
	ldr r0, [sp, r2, lsl #2]
	push {r0}
`, label), nil

	case code.OpCall:
		if len(operands) != 1 {
			return "", fmt.Errorf("OpCall should have only one operand")
		}
		numArgs, ok := operands[0].(int)
		if !ok {
			return "", fmt.Errorf("OpCall argument should be an integer")
		}
		// should jump to the function memory position (on stack before arguments), set the fp
		// and push the returned value (on r0), storing fp and lr for the future
		// the `add r1, pc, #4` is not #12 because to the pc being advanced two positions due
		// to pipelined execution
		return fmt.Sprintf(`%s:  @OpCall
	push {lr}
	add sp, sp, #%d
	ldr r0, [sp]
	str fp, [sp]
	mov fp, sp
	add r1, pc, #4
	mov lr, r1
	mov pc, r0
	push {r0}
`, label, (numArgs+1)*4), nil

	case code.OpReturnValue:
		if len(operands) != 1 {
			return "", fmt.Errorf("OpReturnValue should only one operand")
		}
		// old lr to be restored from memory
		lrIndex, ok := operands[0].(int)
		if !ok {
			return "", fmt.Errorf("OpReturnValue argument should be an integer")
		}
		return fmt.Sprintf(`%s:  @OpReturnValue
	pop {r0}
	mov r1, lr
	ldr lr, [fp, #-%d]
	mov sp, fp
	ldr fp, [sp]
	add sp, sp, #4
	mov pc, r1
`, label, lrIndex*4), nil

	case code.OpReturn:
		if len(operands) != 1 {
			return "", fmt.Errorf("OpReturn should have only one operands")
		}
		// old lr to be restored from memory
		lrIndex, ok := operands[0].(int)
		if !ok {
			return "", fmt.Errorf("OpReturn argument should be an integer")
		}
		// empty returns return null!
		return fmt.Sprintf(`%s:  @OpReturn
	mov r0, #0
	mov r1, lr
	ldr lr, [fp, #-%d]
	mov sp, fp
	ldr fp, [sp]
	add sp, sp, #4
	mov pc, r1
`, label, lrIndex*4), nil

	case code.OpGetLocal:
		if len(operands) != 1 {
			return "", fmt.Errorf("OpGetLocal should have only one operand")
		}
		localIndex, ok := operands[0].(int)
		if !ok {
			return "", fmt.Errorf("OpGetLocal first argument should be an integer")
		}
		// (localIndex+1)*4 because the fp points to the previous sp
		return fmt.Sprintf(`%s:  @OpGetLocal
	sub r0, fp, #%d
	ldr r1, [r0]
	push {r1}
`, label, (localIndex+1)*4), nil

	case code.OpSetLocal:
		if len(operands) != 1 {
			return "", fmt.Errorf("OpSetLocal should have only one operands")
		}
		localIndex, ok := operands[0].(int)
		if !ok {
			return "", fmt.Errorf("OpSetLocal first argument should be an integer")
		}
		// (localIndex+1)*4 because the fp points to the previous sp
		return fmt.Sprintf(`%s:  @OpSetLocal
	sub r0, fp, #%d
	pop {r1}
	str r0, [r1]
`, label, (localIndex+1)*4), nil

	case code.OpGetBuiltin:
		// TODO: implement builtins get and call
		if len(operands) != 1 {
			return "", fmt.Errorf("OpGetBuiltin should have only one operand")
		}
		return fmt.Sprintf(`%s:  @OpGetBuiltin
	ldr r0, =%s
	push {r0}
`, label, operands[0]), nil
	}

	return "", fmt.Errorf("unknown operator: %d", op)
}

func makeBinaryOperation(op code.Opcode, label string) string {
	format := `%s:  @Op%s
	pop {r1, r2}
	%s
	push {r0}
`

	switch op {
	case code.OpAdd:
		return fmt.Sprintf(format, label, "Add", "add r0, r2, r1")
	case code.OpSub:
		return fmt.Sprintf(format, label, "Sub", "sub r0, r2, r1")
	case code.OpMul:
		return fmt.Sprintf(format, label, "Mul", "mul r0, r2, r1")
	case code.OpDiv:
		// TODO: check if args order is right
		// calls ABI integer division routine
		return fmt.Sprintf(format, label, "Div", "bl __aeabi_idiv")
	}

	return ""
}

func makeBoolean(value bool, label string) string {
	format := `%s:  @Op%s
	mov r0, #%d
	push {r0}
`

	if value {
		return fmt.Sprintf(format, label, "True", 1)
	}

	return fmt.Sprintf(format, label, "False", 0)
}

func makeComparison(op code.Opcode, label string) string {
	format := `%s:  @Op%s
	mov r0, #0
	pop {r1, r2}
	cmp r1, r2
	mov%s r0, #1
	push {r0}
`

	switch op {
	case code.OpEqual:
		return fmt.Sprintf(format, label, "Equal", "eq")
	case code.OpNotEqual:
		return fmt.Sprintf(format, label, "NotEqual", "neq")
	case code.OpGreaterThan:
		return fmt.Sprintf(format, label, "GreaterThan", "gt")
	}

	return ""
}

func MakeFunctionPreamble(numOfLocals int) string {
	return fmt.Sprintf("	sub sp, sp, #%d\n", numOfLocals)

}
