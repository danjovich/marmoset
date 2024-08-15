.global _start
.text
_start:
L0:  @OpConstant
	mov r0, #52
	push {r0}

L3:  @OpSetGlobal
	pop {r0}
	ldr r1, =_a
	str r0, [r1]

L6:  @OpConstant
	mov r0, #34
	push {r0}

L9:  @OpSetGlobal
	pop {r0}
	ldr r1, =_b
	str r0, [r1]

L12:  @OpConstant
	ldr r0, =sub
	push {r0}

L15:  @OpSetGlobal
	pop {r0}
	ldr r1, =_sub
	str r0, [r1]

L18:  @OpGetGlobal
	ldr r0, =_sub
	ldr r1, [r0]
	push {r1}

L21:  @OpGetGlobal
	ldr r0, =_a
	ldr r1, [r0]
	push {r1}

L24:  @OpGetGlobal
	ldr r0, =_b
	ldr r1, [r0]
	push {r1}

L27:  @OpCall
	push {lr}
	add sp, sp, #12
	ldr r0, [sp]
	str fp, [sp]
	mov fp, sp
	add r1, pc, #4
	mov lr, r1
	mov pc, r0
	push {r0}

L29:  @OpSetGlobal
	pop {r0}
	ldr r1, =_c
	str r0, [r1]

_end: 
	mov r0, #0 
	mov r7, #1 
	svc #0

sub:
	sub sp, sp, #12

L0_sub:  @OpGetLocal
	sub r0, fp, #4
	ldr r1, [r0]
	push {r1}

L2_sub:  @OpGetLocal
	sub r0, fp, #8
	ldr r1, [r0]
	push {r1}

L4_sub:  @OpSub
	pop {r1, r2}
	sub r0, r2, r1
	push {r0}

L5_sub:  @OpReturnValue
	pop {r0}
	mov r1, lr
	ldr lr, [fp, #-12]
	mov sp, fp
	ldr fp, [sp]
	add sp, sp, #4
	mov pc, r1

_sub: .word sub
_a: .word 0x0
_b: .word 0x0
_c: .word 0x0
