package arm

import (
	"marmoset/code"
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op                      code.Opcode
		index                   int
		label                   string
		args                    []any
		expected                string
		expectedPushesMinusPops int
	}{
		{
			code.OpConstant,
			0,
			"_label",
			[]any{"#1"},
			`L0_label:  @OpConstant
	mov r0, #1
	push {r0}
`,
			1,
		},
		{
			code.OpConstant,
			0,
			"_label",
			[]any{"#1", "#2", "#50"},
			`L0_label:  @OpConstant
	mov r0, #1
	push {r0}
	mov r0, #2
	push {r0}
	mov r0, #50
	push {r0}
`,
			3,
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
			-1,
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
			-1,
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
			-1,
		},
		{
			code.OpDiv,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpDiv
	pop {r1, r2}
	bl __aeabi_idiv
	push {r0}
`,
			-1,
		},
		{
			code.OpPop,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpPop
	pop {r0}
`,
			-1,
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
			1,
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
			1,
		},
		{
			code.OpEqual,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpEqual
	mov r0, #0
	pop {r1, r2}
	cmp r1, r2
	moveq r0, #1
	push {r0}
`,
			-1,
		},
		{
			code.OpNotEqual,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpNotEqual
	mov r0, #0
	pop {r1, r2}
	cmp r1, r2
	movneq r0, #1
	push {r0}
`,
			-1,
		},
		{
			code.OpGreaterThan,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpGreaterThan
	mov r0, #0
	pop {r1, r2}
	cmp r1, r2
	movgt r0, #1
	push {r0}
`,
			-1,
		},
		{
			code.OpMinus,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpMinus
	pop {r0}
	sub r0, #0, r0
	push {r0}
`,
			0,
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
			0,
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
			-1,
		},
		{
			code.OpJump,
			0,
			"_label",
			[]any{"L2"},
			`L0_label:  @OpJump
	b L2
`,
			0,
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
			1,
		},
		{
			code.OpGetGlobal,
			0,
			"_label",
			[]any{"global_var"},
			`L0_label:  @OpGetGlobal
	ldr r0, #_global_var
	push {r0}
`,
			1,
		},
		{
			code.OpSetGlobal,
			0,
			"_label",
			[]any{"global_var"},
			`L0_label:  @OpSetGlobal
	pop {r0}
	str r0, #_global_var
`,
			-1,
		},
		{
			code.OpArray,
			0,
			"_label",
			[]any{3},
			`L0_label:  @OpArray
	mov r0, #3
	push {r0}
	sub r0, sp, #4
	push {r0}
`,
			2,
		},
		{
			code.OpIndex,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpIndex
	pop {r3} 
	ldr r0, [r3, #4]
	ldr r1, [r3, #8]
	sub r2, r1, r0
	sub r2, r2, #1
	ldr r0, [sp, r2, lsl #2]
	push {r0}
`,
			0,
		},
		{
			code.OpCall,
			0,
			"_label",
			[]any{2},
			`L0_label:  @OpCall
	add r0, sp, #8
	mov fp, sp
	ldr r0, [r0]
	blx r0
	push {r0}
`,
			-1,
		},
		{
			code.OpReturnValue,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpReturnValue
	pop {r0}
	mov sp, fp
	mov pc, lr
`,
			0,
		},
		{
			code.OpReturn,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpReturn
	mov sp, fp
	mov r0, #0
	mov pc, lr
`,
			0,
		},
		{
			code.OpGetLocal,
			0,
			"_label",
			[]any{3},
			`L0_label:  @OpGetLocal
	add r0, sp, #12
	ldr r1, [r0]
	push {r1}
`,
			1,
		},
		{
			code.OpSetLocal,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpSetLocal
`,
			0,
		},
		{
			code.OpGetBuiltin,
			0,
			"_label",
			[]any{},
			`L0_label:  @OpGetBuiltin
`,
			0,
		},
	}

	for _, test := range tests {
		result, pushesMinusPops, err := Make(test.op, test.index, test.label, test.args...)

		if err != nil {
			t.Errorf("Make errored: %s", err)
		}

		if result != test.expected {
			t.Errorf("generated assembler source is wrong.\nwant=\n\"%s\"\ngot=\n\"%s\"",
				test.expected, result)
		}

		if pushesMinusPops != test.expectedPushesMinusPops {
			t.Errorf("opcode %v: pushed - pops is wrong.want=%d, got=%d",
				test.op, test.expectedPushesMinusPops, pushesMinusPops)
		}
	}
}
