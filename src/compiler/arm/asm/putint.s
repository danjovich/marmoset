putint:
	sub sp, sp, #8

L0_putint:  @OpConstant
	mov r0, #0
	push {r0}

L3_putint:  @OpGetLocal
	sub r0, fp, #4
	ldr r1, [r0]
	push {r1}

L5_putint:  @OpGreaterThan
	mov r0, #0
	pop {r1, r2}
	cmp r2, r1
	movgt r0, #1
	push {r0}

L6_putint:  @OpJumpNotTruthy
	pop {r0}
	cmp r0, #0
	beq L25_putint

L9_putint:  @OpGetBuiltin
	ldr r0, =put
	push {r0}

L11_putint:  @OpConstant
	mov r0, #45
	push {r0}

L14_putint:  @OpCall
	push {lr}
	add sp, sp, #8
	ldr r0, [sp]
	str fp, [sp]
	mov fp, sp
	add r1, pc, #4
	mov lr, r1
	mov pc, r0
	push {r0}

L16_putint:  @OpPop
	pop {r0}

L17_putint:  @OpGetLocal
	sub r0, fp, #4
	ldr r1, [r0]
	push {r1}

L19_putint:  @OpMinus
	pop {r0}
	mov r1, #0
	sub r0, r1, r0
	push {r0}

L20_putint:  @OpSetLocal
	sub r0, fp, #4
	pop {r1}
	str r1, [r0]

L22_putint:  @OpJump
	b L26_putint

L25_putint:  @OpNull
	mov r0, #0
	push {r0}

L26_putint:  @OpPop
	pop {r0}

L27_putint:  @OpGetLocal
	sub r0, fp, #4
	ldr r1, [r0]
	push {r1}

L29_putint:  @OpConstant
	mov r0, #10
	push {r0}

L32_putint:  @OpDiv
	pop {r1, r2}
	sdiv r0, r2, r1
	push {r0}

L33_putint:  @OpJumpNotTruthy
	pop {r0}
	cmp r0, #0
	beq L50_putint

L36_putint:  @OpGetGlobal
	ldr r0, =_putint
	ldr r1, [r0]
	push {r1}

L39_putint:  @OpGetLocal
	sub r0, fp, #4
	ldr r1, [r0]
	push {r1}

L41_putint:  @OpConstant
	mov r0, #10
	push {r0}

L44_putint:  @OpDiv
	pop {r1, r2}
	sdiv r0, r2, r1
	push {r0}

L45_putint:  @OpCall
	push {lr}
	add sp, sp, #8
	ldr r0, [sp]
	str fp, [sp]
	mov fp, sp
	add r1, pc, #4
	mov lr, r1
	mov pc, r0
	push {r0}

L47_putint:  @OpJump
	b L51_putint

L50_putint:  @OpNull
	mov r0, #0
	push {r0}

L51_putint:  @OpPop
	pop {r0}

L52_putint:  @OpGetBuiltin
	ldr r0, =put
	push {r0}

L54_putint:  @OpGetLocal
	sub r0, fp, #4
	ldr r1, [r0]
	push {r1}

L56_putint:  @OpConstant
	mov r0, #10
	push {r0}

L59_putint:  @OpRest
	pop {r1, r2}
	udiv r0, r2, r1
	mls r0, r0, r1, r2
	push {r0}

L60_putint:  @OpConstant
	mov r0, #48
	push {r0}

L63_putint:  @OpAdd
	pop {r1, r2}
	add r0, r2, r1
	push {r0}

L64_putint:  @OpCall
	push {lr}
	add sp, sp, #8
	ldr r0, [sp]
	str fp, [sp]
	mov fp, sp
	add r1, pc, #4
	mov lr, r1
	mov pc, r0
	push {r0}

L66_putint:  @OpReturnValue
	pop {r0}
	mov r1, lr
	ldr lr, [fp, #-8]
	mov sp, fp
	ldr fp, [sp]
	add sp, sp, #4
	mov pc, r1

_putint: .word putint
