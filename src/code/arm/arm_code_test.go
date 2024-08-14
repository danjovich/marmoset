package arm

import (
	"marmoset/code"
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       code.Opcode
		index    int
		args     []interface{}
		expected string
	}{
		{
			code.OpConstant,
			0,
			[]interface{}{"#1"},
			`L0:  @OpConstant
	mov r0, #1
	push {r0}`,
		},
		{
			code.OpAdd,
			0,
			[]interface{}{},
			`L0:  @OpAdd
	pop {r1, r2}
	add r0, r1, r2
	push {r0}`,
		},
		{
			code.OpSub,
			0,
			[]interface{}{},
			`L0:  @OpSub
	pop {r1, r2}
	sub r0, r1, r2
	push {r0}`,
		},
		{
			code.OpMul,
			0,
			[]interface{}{},
			`L0:  @OpMul
	pop {r1, r2}
	mul r0, r1, r2
	push {r0}`,
		},
		{
			code.OpDiv,
			0,
			[]interface{}{},
			`L0:  @OpDiv
	pop {r1, r2}
	bl __aeabi_idiv
	push {r0}`,
		},
		{
			code.OpPop,
			0,
			[]interface{}{},
			`L0:  @OpPop
	pop {r0}`,
		},
		{
			code.OpTrue,
			0,
			[]interface{}{},
			`L0:  @OpTrue
	mov r0, #1
	push {r0}`,
		},
		{
			code.OpFalse,
			0,
			[]interface{}{},
			`L0:  @OpFalse
	mov r0, #0
	push {r0}`,
		},
		{
			code.OpEqual,
			0,
			[]interface{}{},
			`L0:  @OpEqual
	mov r0, #0
	pop {r1, r2}
	cmp r1, r2
	moveq r0, #1
	push {r0}`,
		},
		{
			code.OpNotEqual,
			0,
			[]interface{}{},
			`L0:  @OpNotEqual
	mov r0, #0
	pop {r1, r2}
	cmp r1, r2
	movneq r0, #1
	push {r0}`,
		},
		{
			code.OpGreaterThan,
			0,
			[]interface{}{},
			`L0:  @OpGreaterThan
	mov r0, #0
	pop {r1, r2}
	cmp r1, r2
	movgt r0, #1
	push {r0}`,
		},
		{
			code.OpMinus,
			0,
			[]interface{}{},
			`L0:  @OpMinus
	pop {r0}
	sub r0, #0, r0
	push {r0}`,
		},
		{
			code.OpBang,
			0,
			[]interface{}{},
			`L0:  @OpBang
	mov r1, #0
	pop {r0}
	cmp r0, #0
	moveq r1, #1
	push {r1}`,
		},
		{
			code.OpJumpNotTruthy,
			0,
			[]interface{}{"L2"},
			`L0:  @OpJumpNotTruthy
	pop {r0}
	cmp r0, #0
	beq L2`,
		},
		{
			code.OpJump,
			0,
			[]interface{}{"L2"},
			`L0:  @OpJump
	b L2`,
		},
		{
			code.OpNull,
			0,
			[]interface{}{},
			`L0:  @OpNull
	mov r0, #0
	push {r0}`,
		},
		{
			code.OpGetGlobal,
			0,
			[]interface{}{"global_var"},
			`L0:  @OpGetGlobal
	ldr r0, =global_var
	push {r0}`,
		},
		{
			code.OpSetGlobal,
			0,
			[]interface{}{"global_var"},
			`L0:  @OpSetGlobal
	pop {r0}
	str r0, =global_var`,
		},
		{
			code.OpArray,
			0,
			[]interface{}{3},
			`L0:  @OpArray
	mov r0, #3
	push {r0}`,
		},
		{
			code.OpIndex,
			0,
			[]interface{}{},
			`L0:  @OpIndex
	pop {r0, r1}
	sub r2, r1, r0
	sub r2, r2, #1
	ldr r0, [sp, r2, lsl #2]
	push {r0}
	`,
		},
		{
			code.OpCall,
			0,
			[]interface{}{2},
			`L0:  @OpCall
	add r0, sp, #8
	mov fp, r0
	ldr pc, [fp]`,
		},
		{
			code.OpReturnValue,
			0,
			[]interface{}{},
			`L0:  @OpReturnValue
	pop {r0}
	mov sp, fp
	push {r0}
	b lr`,
		},
		{
			code.OpReturn,
			0,
			[]interface{}{},
			`L0:  @OpReturn
	mov sp, fp
	mov r0, #0
	push {r0}
	b lr`,
		},
		{
			code.OpGetLocal,
			0,
			[]interface{}{"local_var"},
			`L0:  @OpGetLocal
	ldr r0, =local_var
	push {r0}`,
		},
		{
			code.OpSetLocal,
			0,
			[]interface{}{"local_var"},
			`L0:  @OpSetLocal
	pop {r0}
	str r0, =local_var`,
		},
		{
			code.OpGetBuiltin,
			0,
			[]interface{}{},
			`L0:  @OpGetBuiltin`,
		},
	}

	for _, test := range tests {
		result, err := Make(test.op, test.index, test.args...)

		if err != nil {
			t.Errorf("Make errored: %s", err)
		}

		if result != test.expected {
			t.Errorf("generated assembler source is wrong.\nwant=\n\"%s\"\ngot=\n\"%s\"",
				test.expected, result)
		}
	}
}
