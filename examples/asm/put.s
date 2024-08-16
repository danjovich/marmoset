.global _start
.text
_start:
L0:  @OpConstant
	ldr r0, =putABC
	push {r0}

L3:  @OpSetGlobal
	pop {r0}
	ldr r1, =_putABC
	str r0, [r1]

L6:  @OpGetGlobal
	ldr r0, =_putABC
	ldr r1, [r0]
	push {r1}

L9:  @OpCall
	push {lr}
	add sp, sp, #4
	ldr r0, [sp]
	str fp, [sp]
	mov fp, sp
	add r1, pc, #4
	mov lr, r1
	mov pc, r0
	push {r0}

L11:  @OpPop
	pop {r0}

_end: 
	mov r0, #0 
	mov r7, #1 
	svc #0

putABC:
	sub sp, sp, #4

L0_putABC:  @OpGetBuiltin
	ldr r0, =put
	push {r0}

L2_putABC:  @OpConstant
	mov r0, #65
	push {r0}

L5_putABC:  @OpCall
	push {lr}
	add sp, sp, #8
	ldr r0, [sp]
	str fp, [sp]
	mov fp, sp
	add r1, pc, #4
	mov lr, r1
	mov pc, r0
	push {r0}

L7_putABC:  @OpPop
	pop {r0}

L8_putABC:  @OpGetBuiltin
	ldr r0, =put
	push {r0}

L10_putABC:  @OpConstant
	mov r0, #66
	push {r0}

L13_putABC:  @OpCall
	push {lr}
	add sp, sp, #8
	ldr r0, [sp]
	str fp, [sp]
	mov fp, sp
	add r1, pc, #4
	mov lr, r1
	mov pc, r0
	push {r0}

L15_putABC:  @OpPop
	pop {r0}

L16_putABC:  @OpGetBuiltin
	ldr r0, =put
	push {r0}

L18_putABC:  @OpConstant
	mov r0, #67
	push {r0}

L21_putABC:  @OpCall
	push {lr}
	add sp, sp, #8
	ldr r0, [sp]
	str fp, [sp]
	mov fp, sp
	add r1, pc, #4
	mov lr, r1
	mov pc, r0
	push {r0}

L23_putABC:  @OpReturnValue
	pop {r0}
	mov r1, lr
	ldr lr, [fp, #-4]
	mov sp, fp
	ldr fp, [sp]
	add sp, sp, #4
	mov pc, r1

put:
	sub sp, sp, #4

L0_put: @put
	mov r0, #1
	add r1, fp, #-4
	mov r2, #1 
	mov r7, #4 
	svc #0

L3put:  @OpReturn
	mov r0, #0
	mov r1, lr
	ldr lr, [fp, #-8]
	mov sp, fp
	ldr fp, [sp]
	add sp, sp, #4
	mov pc, r1

_putABC: .word putABC
