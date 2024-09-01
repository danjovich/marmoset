# machine name
UNAME_M := $(shell uname -m)

# checks if the machine is NOT an ARM
ifeq ($(filter arm%,$(UNAME_M)),)
# adds prefix to binutils calls
	BINUTILS_PREFIX := arm-linux-gnueabihf-
endif

# assembler command
AS := $(BINUTILS_PREFIX)as
# linker command
LD := $(BINUTILS_PREFIX)ld

# buids the Marmoset compiler
bin/marmoset : src/*.go src/**/*.go src/**/**/*.go
	cd src && go build -o ../$@ main.go

# cleans Go module cache
clean :
	cd src && go clean -modcache

# builds the assembly examples using the Marmoset compiler
examples/asm/%.s : bin/marmoset examples/%.marm
	if [ ! -d examples/asm ]; then mkdir examples/asm; fi
	./$^ > $@

# builds the binary from an assembler source
examples/bin/%.out : examples/asm/%.s
	if [ ! -d examples/bin ]; then mkdir examples/bin; fi
	$(AS) -o /tmp/a.out -g $?
# ideally, -N shouldn't be used in ld, as it makes all code rwx,
# but this made it easier to store global data (since there is no 
# .data section)
	$(LD) -N -o $@ -g /tmp/a.out
	rm /tmp/a.out

# runs qemu for a certain example (unnecessary on ARM machines)
qemu-% : examples/bin/%.out
	qemu-arm -L /usr/arm-linux-gnueabihf -g 1234 $?

# runs gdb for a certain example (qemu must be running if not on ARM)
gdb-% : examples/bin/%.out
	gdb-multiarch -q --nh \
		-ex 'set architecture arm' \
		-ex 'file $?' \
		-ex 'target remote localhost:1234' \
		-ex 'layout split' \
		-ex 'layout regs'

# runs an example's elf file
run-% : examples/bin/%.out
	./$?
