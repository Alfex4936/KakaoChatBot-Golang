build:
	go build -o bin/main main/main.go

run:
	go run main.go

run-linux:
	./bin/main

run-windows:
	./bin/main-windows-amd64.exe

compile:
	# 32-Bit Systems
	# FreeBDS
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 main.go
	# MacOS
	GOOS=darwin GOARCH=386 go build -o bin/main-darwin-386 main.go
	# Linux
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 main.go
	# Windows
	GOOS=windows GOARCH=386 go build -o bin/main-windows-386 main.go
	# 64-Bit
	# FreeBDS
	GOOS=freebsd GOARCH=amd64 go build -o bin/main-freebsd-amd64 main.go
	# MacOS
	GOOS=darwin GOARCH=amd64 go build -o bin/main-darwin-amd64 main.go
	# Linux
	GOOS=linux GOARCH=amd64 go build -o bin/main-linux-amd64 main.go
	# Windows
	GOOS=windows GOARCH=amd64 go build -o bin/main-windows-amd64 main.go

compile-windows:
	@echo "Compiling for every OS and Platform"
	set GOOS=freebsd& set GOARCH=amd64& go build -o bin/main-freebsd-amd64 main/main.go
	set GOOS=darwin& set GOARCH=amd64& go build -o bin/main-darwin-amd64 main/main.go
	set GOOS=linux& set GOARCH=amd64& go build -o bin/main-linux-amd64 main/main.go
	set GOOS=windows& set GOARCH=amd64& go build -o bin/main-windows-amd64.exe main/main.go

all: build run