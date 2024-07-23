ifeq ($(OS),Windows_NT)
	EXECUTABLE_EXTENSION := .exe
endif

bin/monkey$(EXECUTABLE_EXTENSION) : src/*.go src/**/*.go
	cd src && go build -o ../$@ main.go

run : bin/monkey$(EXECUTABLE_EXTENSION)
	./$?