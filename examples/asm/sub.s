.global _start
_start:
	mov sp, #0x4000
	mov fp, sp

L0:  @OpConstant
	mov r0, #52
	push {r0}

L3:  @OpSetGlobal
	pop {r0}
	str r0, #_a

L6:  @OpConstant
	mov r0, #34
	push {r0}

L9:  @OpSetGlobal
	pop {r0}
	str r0, #_b

L12:  @OpConstant
	mov r0, #sub
	push {r0}

L15:  @OpSetGlobal
	pop {r0}
	str r0, #_sub

L18:  @OpGetGlobal
	ldr r0, #_sub
	push {r0}

L21:  @OpGetGlobal
	ldr r0, #_a
	push {r0}

L24:  @OpGetGlobal
	ldr r0, #_b
	push {r0}

L27:  @OpCall
	push {lr}
	add sp, sp, #12
	ldr r0, [sp]
	str fp, [sp]
	mov fp, sp
	blx r0
	push {r0}

L29:  @OpSetGlobal
	pop {r0}
	str r0, #_c

_end: b _end

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

_c: .word 0x0
_sub: .word sub
_a: .word 0x0
_b: .word 0x0
