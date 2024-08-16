putintln:
	sub sp, sp, #8

L0_putintln:  @OpGetBuiltin
	ldr r0, =putint
	push {r0}

L2_putintln:  @OpGetLocal
	sub r0, fp, #4
	ldr r1, [r0]
	push {r1}

L4_putintln:  @OpCall
	push {lr}
	add sp, sp, #8
	ldr r0, [sp]
	str fp, [sp]
	mov fp, sp
	add r1, pc, #4
	mov lr, r1
	mov pc, r0
	push {r0}

L6_putintln:  @OpPop
	pop {r0}

L7_putintln:  @OpGetBuiltin
	ldr r0, =put
	push {r0}

L9_putintln:  @OpConstant
	mov r0, #10
	push {r0}

L12_putintln:  @OpCall
	push {lr}
	add sp, sp, #8
	ldr r0, [sp]
	str fp, [sp]
	mov fp, sp
	add r1, pc, #4
	mov lr, r1
	mov pc, r0
	push {r0}

L14_putintln:  @OpReturnValue
	pop {r0}
	mov r1, lr
	ldr lr, [fp, #-8]
	mov sp, fp
	ldr fp, [sp]
	add sp, sp, #4
	mov pc, r1

_putintln: .word putintln
