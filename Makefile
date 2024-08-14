ifeq ($(OS),Windows_NT)
	EXECUTABLE_EXTENSION := .exe
endif

bin/marmoset$(EXECUTABLE_EXTENSION) : src/*.go src/**/*.go src/**/**/*.go
	cd src && go build -o ../$@ main.go

run : bin/marmoset$(EXECUTABLE_EXTENSION)
	./$?

clean :
	cd src && go clean -modcache