.global _start
_start:
	mov sp, #0x4000
	mov fp, sp

L0:  @OpConstant
	mov r0, #fibonacci
	push {r0}

L3:  @OpSetGlobal
	pop {r0}
	str r0, #_fibonacci

L6:  @OpGetGlobal
	ldr r0, #_fibonacci
	push {r0}

L9:  @OpConstant
	mov r0, #15
	push {r0}

L12:  @OpCall
	push {lr}
	add sp, sp, #8
	ldr r0, [sp]
	str fp, [sp]
	mov fp, sp
	blx r0
	push {r0}

L14:  @OpPop
	pop {r0}

_end: b _end

fibonacci:
	sub sp, sp, #8

L0_fibonacci:  @OpGetLocal
	sub r0, fp, #4
	ldr r1, [r0]
	push {r1}

L2_fibonacci:  @OpConstant
	mov r0, #0
	push {r0}

L5_fibonacci:  @OpEqual
	mov r0, #0
	pop {r1, r2}
	cmp r1, r2
	moveq r0, #1
	push {r0}

L6_fibonacci:  @OpJumpNotTruthy
	pop {r0}
	cmp r0, #0
	beq L16_fibonacci

L9_fibonacci:  @OpConstant
	mov r0, #0
	push {r0}

L12_fibonacci:  @OpReturnValue
	pop {r0}
	mov r1, lr
	ldr lr, [fp, #-8]
	mov sp, fp
	ldr fp, [sp]
	add sp, sp, #4
	mov pc, r1

L13_fibonacci:  @OpJump
	b L56_fibonacci

L16_fibonacci:  @OpGetLocal
	sub r0, fp, #4
	ldr r1, [r0]
	push {r1}

L18_fibonacci:  @OpConstant
	mov r0, #1
	push {r0}

L21_fibonacci:  @OpEqual
	mov r0, #0
	pop {r1, r2}
	cmp r1, r2
	moveq r0, #1
	push {r0}

L22_fibonacci:  @OpJumpNotTruthy
	pop {r0}
	cmp r0, #0
	beq L32_fibonacci

L25_fibonacci:  @OpConstant
	mov r0, #1
	push {r0}

L28_fibonacci:  @OpReturnValue
	pop {r0}
	mov r1, lr
	ldr lr, [fp, #-8]
	mov sp, fp
	ldr fp, [sp]
	add sp, sp, #4
	mov pc, r1

L29_fibonacci:  @OpJump
	b L56_fibonacci

L32_fibonacci:  @OpGetGlobal
	ldr r0, #_fibonacci
	push {r0}

L35_fibonacci:  @OpGetLocal
	sub r0, fp, #4
	ldr r1, [r0]
	push {r1}

L37_fibonacci:  @OpConstant
	mov r0, #1
	push {r0}

L40_fibonacci:  @OpSub
	pop {r1, r2}
	sub r0, r2, r1
	push {r0}

L41_fibonacci:  @OpCall
	push {lr}
	add sp, sp, #8
	ldr r0, [sp]
	str fp, [sp]
	mov fp, sp
	blx r0
	push {r0}

L43_fibonacci:  @OpGetGlobal
	ldr r0, #_fibonacci
	push {r0}

L46_fibonacci:  @OpGetLocal
	sub r0, fp, #4
	ldr r1, [r0]
	push {r1}

L48_fibonacci:  @OpConstant
	mov r0, #2
	push {r0}

L51_fibonacci:  @OpSub
	pop {r1, r2}
	sub r0, r2, r1
	push {r0}

L52_fibonacci:  @OpCall
	push {lr}
	add sp, sp, #8
	ldr r0, [sp]
	str fp, [sp]
	mov fp, sp
	blx r0
	push {r0}

L54_fibonacci:  @OpAdd
	pop {r1, r2}
	add r0, r2, r1
	push {r0}

L55_fibonacci:  @OpReturnValue
	pop {r0}
	mov r1, lr
	ldr lr, [fp, #-8]
	mov sp, fp
	ldr fp, [sp]
	add sp, sp, #4
	mov pc, r1

L56_fibonacci:  @OpReturnValue
	pop {r0}
	mov r1, lr
	ldr lr, [fp, #-8]
	mov sp, fp
	ldr fp, [sp]
	add sp, sp, #4
	mov pc, r1

_fibonacci: .word fibonacci
