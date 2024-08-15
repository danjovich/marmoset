.global _start
.text
_start:
L0:  @OpConstant
	ldr r0, =factorial
	push {r0}

L3:  @OpSetGlobal
	pop {r0}
	ldr r1, =_factorial
	str r0, [r1]

L6:  @OpGetGlobal
	ldr r0, =_factorial
	ldr r1, [r0]
	push {r1}

L9:  @OpConstant
	mov r0, #10
	push {r0}

L12:  @OpCall
	push {lr}
	add sp, sp, #8
	ldr r0, [sp]
	str fp, [sp]
	mov fp, sp
	add r1, pc, #4
	mov lr, r1
	mov pc, r0
	push {r0}

L14:  @OpPop
	pop {r0}

_end: 
	mov r0, #0 
	mov r7, #1 
	svc #0

factorial:
	sub sp, sp, #8

L0_factorial:  @OpGetLocal
	sub r0, fp, #4
	ldr r1, [r0]
	push {r1}

L2_factorial:  @OpConstant
	mov r0, #0
	push {r0}

L5_factorial:  @OpEqual
	mov r0, #0
	pop {r1, r2}
	cmp r1, r2
	moveq r0, #1
	push {r0}

L6_factorial:  @OpJumpNotTruthy
	pop {r0}
	cmp r0, #0
	beq L16_factorial

L9_factorial:  @OpConstant
	mov r0, #1
	push {r0}

L12_factorial:  @OpReturnValue
	pop {r0}
	mov r1, lr
	ldr lr, [fp, #-8]
	mov sp, fp
	ldr fp, [sp]
	add sp, sp, #4
	mov pc, r1

L13_factorial:  @OpJump
	b L17_factorial

L16_factorial:  @OpNull
	mov r0, #0
	push {r0}

L17_factorial:  @OpPop
	pop {r0}

L18_factorial:  @OpGetGlobal
	ldr r0, =_factorial
	ldr r1, [r0]
	push {r1}

L21_factorial:  @OpGetLocal
	sub r0, fp, #4
	ldr r1, [r0]
	push {r1}

L23_factorial:  @OpConstant
	mov r0, #1
	push {r0}

L26_factorial:  @OpSub
	pop {r1, r2}
	sub r0, r2, r1
	push {r0}

L27_factorial:  @OpCall
	push {lr}
	add sp, sp, #8
	ldr r0, [sp]
	str fp, [sp]
	mov fp, sp
	add r1, pc, #4
	mov lr, r1
	mov pc, r0
	push {r0}

L29_factorial:  @OpGetLocal
	sub r0, fp, #4
	ldr r1, [r0]
	push {r1}

L31_factorial:  @OpMul
	pop {r1, r2}
	mul r0, r2, r1
	push {r0}

L32_factorial:  @OpReturnValue
	pop {r0}
	mov r1, lr
	ldr lr, [fp, #-8]
	mov sp, fp
	ldr fp, [sp]
	add sp, sp, #4
	mov pc, r1

_factorial: .word factorial
