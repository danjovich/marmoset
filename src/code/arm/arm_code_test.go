package arm

import (
	"marmoset/code"
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       code.Opcode
		index    int
		label    string
		args     []any
		expected string
	}{
		{
			code.OpConstant,
			0,
			"_label",
			[]any{"mov r0, #1"},
			`L0_label:  @OpConstant
	mov r0, #1
	push {r0}
`,
		},
		{
			code.OpConstant,
			0,
			"_label",
			[]any{"mov r0, #1", "mov r0, #2", "mov r0, #50"},
			`L0_label:  @OpConstant
	mov r0, #1
	push {r0}
	mov r0, #2
	push {r0}
	mov r0, #50
	push {r0}
`,
		},
		{
			code.OpAdd,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpAdd
	pop {r1, r2}
	add r0, r2, r1
	push {r0}
`,
		},
		{
			code.OpSub,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpSub
	pop {r1, r2}
	sub r0, r2, r1
	push {r0}
`,
		},
		{
			code.OpMul,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpMul
	pop {r1, r2}
	mul r0, r2, r1
	push {r0}
`,
		},
		{
			code.OpDiv,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpDiv
	pop {r1, r2}
	sdiv r0, r2, r1
	push {r0}
`,
		},
		{
			code.OpRest,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpRest
	pop {r1, r2}
	udiv r0, r2, r1
	mls r0, r0, r1, r2
	push {r0}
`,
		},
		{
			code.OpPop,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpPop
	pop {r0}
`,
		},
		{
			code.OpTrue,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpTrue
	mov r0, #1
	push {r0}
`,
		},
		{
			code.OpFalse,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpFalse
	mov r0, #0
	push {r0}
`,
		},
		{
			code.OpEqual,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpEqual
	mov r0, #0
	pop {r1, r2}
	cmp r2, r1
	moveq r0, #1
	push {r0}
`,
		},
		{
			code.OpNotEqual,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpNotEqual
	mov r0, #0
	pop {r1, r2}
	cmp r2, r1
	movne r0, #1
	push {r0}
`,
		},
		{
			code.OpGreaterThan,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpGreaterThan
	mov r0, #0
	pop {r1, r2}
	cmp r2, r1
	movgt r0, #1
	push {r0}
`,
		},
		{
			code.OpMinus,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpMinus
	pop {r0}
	mov r1, #0
	sub r0, r1, r0
	push {r0}
`,
		},
		{
			code.OpBang,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpBang
	mov r1, #0
	pop {r0}
	cmp r0, #0
	moveq r1, #1
	push {r1}
`,
		},
		{
			code.OpJumpNotTruthy,
			0,
			"_label",
			[]any{"L2"},
			`L0_label:  @OpJumpNotTruthy
	pop {r0}
	cmp r0, #0
	beq L2
`,
		},
		{
			code.OpJump,
			0,
			"_label",
			[]any{"L2"},
			`L0_label:  @OpJump
	b L2
`,
		},
		{
			code.OpNull,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpNull
	mov r0, #0
	push {r0}
`,
		},
		{
			code.OpGetGlobal,
			0,
			"_label",
			[]any{"global_var"},
			`L0_label:  @OpGetGlobal
	ldr r0, =_global_var
	ldr r1, [r0]
	push {r1}
`,
		},
		{
			code.OpSetGlobal,
			0,
			"_label",
			[]any{"global_var"},
			`L0_label:  @OpSetGlobal
	pop {r0}
	ldr r1, =_global_var
	str r0, [r1]
`,
		},
		// 		{
		// 			code.OpArray,
		// 			0,
		// 			"_label",
		// 			[]any{3},
		// 			`L0_label:  @OpArray
		// 	mov r0, #3
		// 	push {r0}
		// 	sub r0, sp, #4
		// 	push {r0}
		// `,
		// 		},
		// 		{
		// 			code.OpIndex,
		// 			0,
		// 			"_label",
		// 			[]any{},
		// 			`L0_label:  @OpIndex
		// 	pop {r3}
		// 	ldr r0, [r3, #4]
		// 	ldr r1, [r3, #8]
		// 	sub r2, r1, r0
		// 	sub r2, r2, #1
		// 	ldr r0, [sp, r2, lsl #2]
		// 	push {r0}
		// `,
		// 		},
		{
			code.OpCall,
			0,
			"_label",
			[]any{2},
			`L0_label:  @OpCall
	push {lr}
	add sp, sp, #12
	ldr r0, [sp]
	str fp, [sp]
	mov fp, sp
	add r1, pc, #4
	mov lr, r1
	mov pc, r0
	push {r0}
`,
		},
		{
			code.OpReturnValue,
			0,
			"_label",
			[]any{2},
			`L0_label:  @OpReturnValue
	pop {r0}
	mov r1, lr
	ldr lr, [fp, #-8]
	mov sp, fp
	ldr fp, [sp]
	add sp, sp, #4
	mov pc, r1
`,
		},
		{
			code.OpReturn,
			0,
			"_label",
			[]any{3},
			`L0_label:  @OpReturn
	mov r0, #0
	mov r1, lr
	ldr lr, [fp, #-12]
	mov sp, fp
	ldr fp, [sp]
	add sp, sp, #4
	mov pc, r1
`,
		},
		{
			code.OpGetLocal,
			0,
			"_label",
			[]any{3},
			`L0_label:  @OpGetLocal
	sub r0, fp, #16
	ldr r1, [r0]
	push {r1}
`,
		},
		{
			code.OpSetLocal,
			0,
			"_label",
			[]any{2},
			`L0_label:  @OpSetLocal
	sub r0, fp, #12
	pop {r1}
	str r1, [r0]
`,
		},
		{
			code.OpGetBuiltin,
			0,
			"_label",
			[]any{"put"},
			`L0_label:  @OpGetBuiltin
	ldr r0, =put
	push {r0}
`,
		},
	}

	for _, test := range tests {
		result, err := Make(test.op, test.index, test.label, test.args...)

		if err != nil {
			t.Errorf("Make errored: %s", err)
		}

		if result != test.expected {
			t.Errorf("generated assembler source is wrong.\nwant=\n\"%s\"\ngot=\n\"%s\"",
				test.expected, result)
		}
	}
}
