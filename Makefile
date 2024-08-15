ifeq ($(OS),Windows_NT)
	EXECUTABLE_EXTENSION := .exe
endif

# buids the Marmoset compiler
bin/marmoset$(EXECUTABLE_EXTENSION) : src/*.go src/**/*.go src/**/**/*.go
	cd src && go build -o ../$@ main.go

# cleans Go module cache
clean :
	cd src && go clean -modcache

# builds the assembly examples using the Marmoset compiler
examples/asm/%.s : bin/marmoset$(EXECUTABLE_EXTENSION) examples/%.marm
	./$^ > $@

# builds the binary from an assembler source
examples/bin/%.out : examples/asm/%.s
	if [ ! -d examples/bin ]; then mkdir examples/bin; fi
	arm-linux-gnueabihf-as -o /tmp/a.out -g $?
# ideally, -N shouldn't be used in ld, as it makes all code rwx,
# but this made it easier to store global data (since there is no 
# .data section)
	arm-linux-gnueabihf-ld -N -o $@ -g /tmp/a.out
	rm /tmp/a.out

# runs qemu for a certain example
qemu-% : examples/bin/%.out
	qemu-arm -L /usr/arm-linux-gnueabihf -g 1234 $?

# runs gdb for a certain example (qemu must be running)
gdb-% : examples/bin/%.out
	gdb-multiarch -q --nh \
		-ex 'set architecture arm' \
		-ex 'file $?' \
		-ex 'target remote localhost:1234' \
		-ex 'layout split' \
		-ex 'layout regs'
	

