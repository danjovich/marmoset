package arm

import (
	"fmt"
	"marmoset/code"
)

func Make(op code.Opcode, index int, operands ...interface{}) (string, error) {
	switch op {
	case code.OpConstant:
		if len(operands) != 1 {
			return "", fmt.Errorf("OpConstant should have only one operand")
		}
		return fmt.Sprintf(`L%d:  @OpConstant
	mov r0, %s
	push {r0}`, index, operands[0]), nil

	case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
		if len(operands) != 0 {
			return "", fmt.Errorf("binary operations should not have any operands")
		}
		return makeBinaryOperation(op, index), nil

	case code.OpPop:
		if len(operands) != 0 {
			return "", fmt.Errorf("OpPop should not have any operands")
		}
		return fmt.Sprintf(`L%d:  @OpPop
	pop {r0}`, index), nil

	case code.OpTrue, code.OpFalse:
		if len(operands) != 0 {
			return "", fmt.Errorf("booleans should not have any operands")
		}
		return makeBoolean(op == code.OpTrue, index), nil

	case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
		if len(operands) != 0 {
			return "", fmt.Errorf("comparisons should not have any operands")
		}
		return makeComparison(op, index), nil

	case code.OpMinus:
		if len(operands) != 0 {
			return "", fmt.Errorf("OpMinus should not have any operands")
		}
		return fmt.Sprintf(`L%d:  @OpMinus
	pop {r0}
	sub r0, #0, r0
	push {r0}`, index), nil

	case code.OpBang:
		if len(operands) != 0 {
			return "", fmt.Errorf("OpBang should not have any operands")
		}
		return fmt.Sprintf(`L%d:  @OpBang
	mov r1, #0
	pop {r0}
	cmp r0, #0
	moveq r1, #1
	push {r1}`, index), nil

	case code.OpJumpNotTruthy:
		if len(operands) != 1 {
			return "", fmt.Errorf("OpJumpNotTruthy should have only one operand")
		}
		dest := operands[0]
		return fmt.Sprintf(`L%d:  @OpJumpNotTruthy
	pop {r0}
	cmp r0, #0
	beq %s`, index, dest), nil

	case code.OpJump:
		if len(operands) != 1 {
			return "", fmt.Errorf("OpJump should have only one operand")
		}
		dest := operands[0]
		return fmt.Sprintf(`L%d:  @OpJump
	b %s`, index, dest), nil

	case code.OpNull:
		if len(operands) != 0 {
			return "", fmt.Errorf("OpNull should not have any operands")
		}
		return fmt.Sprintf(`L%d:  @OpNull
	mov r0, #0
	push {r0}`, index), nil

	case code.OpGetGlobal:
		if len(operands) != 1 {
			return "", fmt.Errorf("OpGetGlobal should have only one operand")
		}
		name := operands[0]
		return fmt.Sprintf(`L%d:  @OpGetGlobal
	ldr r0, =%s
	push {r0}`, index, name), nil

	case code.OpSetGlobal:
		if len(operands) != 1 {
			return "", fmt.Errorf("OpSetGlobal should have only one operand")
		}
		name := operands[0]
		return fmt.Sprintf(`L%d:  @OpSetGlobal
	pop {r0}
	str r0, =%s`, index, name), nil

	case code.OpArray:
		if len(operands) != 1 {
			return "", fmt.Errorf("OpArray should have only one operand")
		}
		length := operands[0]
		return fmt.Sprintf(`L%d:  @OpArray
	mov r0, #%d
	push {r0}`, index, length), nil

	case code.OpIndex:
		if len(operands) != 0 {
			return "", fmt.Errorf("OpIndex should not have any operands")
		}
		return fmt.Sprintf(`L%d:  @OpIndex
	pop {r0, r1}
	sub r2, r1, r0
	sub r2, r2, #1
	ldr r0, [sp, r2, lsl #2]
	push {r0}
	`, index), nil

	case code.OpCall:
		if len(operands) != 1 {
			return "", fmt.Errorf("OpCall should have only one operand")
		}
		name := operands[0]
		return fmt.Sprintf(`L%d:  @OpCall
	bl %s`, index, name), nil

	case code.OpReturnValue:
		if len(operands) != 0 {
			return "", fmt.Errorf("OpReturnValue should not have any operands")
		}
		return fmt.Sprintf(`L%d:  @OpReturnValue
	b lr`, index), nil

	case code.OpReturn:
		if len(operands) != 0 {
			return "", fmt.Errorf("OpReturn should not have any operands")
		}
		// empty returns return null!
		return fmt.Sprintf(`L%d:  @OpReturn
	mov r0, #0
	push {r0}
	b lr`, index), nil

	case code.OpGetLocal:
		// TODO: in the compiler, the name here should be name + 'some_hash' to avoid conflict with globals
		if len(operands) != 1 {
			return "", fmt.Errorf("OpGetLocal should have only one operand")
		}
		name := operands[0]
		return fmt.Sprintf(`L%d:  @OpGetLocal
	ldr r0, =%s
	push {r0}`, index, name), nil

	case code.OpSetLocal:
		// TODO: in the compiler, the name here should be name + 'some_hash' to avoid conflict with globals
		if len(operands) != 1 {
			return "", fmt.Errorf("OpSetLocal should have only one operand")
		}
		name := operands[0]
		return fmt.Sprintf(`L%d:  @OpSetLocal
	pop {r0}
	str r0, =%s`, index, name), nil

	case code.OpGetBuiltin:
		// TODO: implement builtins get and call
		if len(operands) != 0 {
			return "", fmt.Errorf("OpGetBuiltin should not have any operands")
		}
		return fmt.Sprintf(`L%d:  @OpGetBuiltin`, index), nil
	}

	return "", nil
	// return "", fmt.Errorf("unknown operator: %d", op)
}

func makeBinaryOperation(op code.Opcode, index int) string {
	format := `L%d:  @Op%s
	pop {r1, r2}
	%s
	push {r0}`

	switch op {
	case code.OpAdd:
		return fmt.Sprintf(format, index, "Add", "add r0, r1, r2")
	case code.OpSub:
		return fmt.Sprintf(format, index, "Sub", "sub r0, r1, r2")
	case code.OpMul:
		return fmt.Sprintf(format, index, "Mul", "mul r0, r1, r2")
	case code.OpDiv:
		// calls ABI integer division routine
		return fmt.Sprintf(format, index, "Div", "bl __aeabi_idiv")
	}

	return ""
}

func makeBoolean(value bool, index int) string {
	format := `L%d:  @Op%s
	mov r0, #%d
	push {r0}`

	if value {
		return fmt.Sprintf(format, index, "True", 1)
	}

	return fmt.Sprintf(format, index, "False", 0)
}

func makeComparison(op code.Opcode, index int) string {
	format := `L%d:  @Op%s
	mov r0, #0
	pop {r1, r2}
	cmp r1, r2
	mov%s r0, #1
	push {r0}`

	switch op {
	case code.OpEqual:
		return fmt.Sprintf(format, index, "Equal", "eq")
	case code.OpNotEqual:
		return fmt.Sprintf(format, index, "NotEqual", "neq")
	case code.OpGreaterThan:
		return fmt.Sprintf(format, index, "GreaterThan", "gt")
	}

	return ""
}
