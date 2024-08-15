ifeq ($(OS),Windows_NT)
	EXECUTABLE_EXTENSION := .exe
endif

# buids the Marmoset compiler
bin/marmoset$(EXECUTABLE_EXTENSION) : src/*.go src/**/*.go src/**/**/*.go
	cd src && go build -o ../$@ main.go

# cleans Go module cache
clean :
	cd src && go clean -modcache

# builds the examples using the Marmoset compiler
examples/asm/%.s : bin/marmoset$(EXECUTABLE_EXTENSION) examples/%.marm
	./$^ > $@
