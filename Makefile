.PHONY:

go:
	go run main/main.go

build:
	go build -ldflags="-w -s" -o bin/main main/main.go
	@echo Upxing...
	upx bin/main
ifneq ($(OS), Windows_NT)
	chmod +x bin/main
endif

run:
ifeq ($(OS), Windows_NT)
# Powershell >= 7
# Execute main without extension ".exe"
	pwsh -NoProfile -ExecutionPolicy Unrestricted -Command "& {Start-Process -FilePath ""./bin/main"" -Wait -NoNewWindow}"
else
	chmod +x ./bin/main
	./bin/main
endif
# go run main/main.go

compile:
ifeq ($(OS), Windows_NT)
	@echo "$(OS): Compiling for every OS and Platform"
	set GOOS=freebsd& set GOARCH=amd64& go build -o bin/main-freebsd-amd64 main/main.go
	set GOOS=darwin& set GOARCH=amd64& go build -o bin/main-darwin-amd64 main/main.go
	set GOOS=linux& set GOARCH=amd64& go build -o bin/main-linux-amd64 main/main.go
	set GOOS=windows& set GOARCH=amd64& go build -o bin/main-windows-amd64.exe main/main.go
else
	@echo "$(OS): Compiling for every OS and Platform"
# 32-Bit Systems
# FreeBDS
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 main/main.go
# MacOS
#	GOOS=darwin GOARCH=386 go build -o bin/main-darwin-386 main/main.go
# Linux
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 main/main.go
# Windows
	GOOS=windows GOARCH=386 go build -o bin/main-windows-386.exe main/main.go
# 64-Bit
# FreeBDS
	GOOS=freebsd GOARCH=amd64 go build -o bin/main-freebsd-amd64 main/main.go
# MacOS
	GOOS=darwin GOARCH=amd64 go build -o bin/main-darwin-amd64 main/main.go
# Linux
	GOOS=linux GOARCH=amd64 go build -o bin/main-linux-amd64 main/main.go
# Windows
	GOOS=windows set GOARCH=amd64 go build -o bin/main-windows-amd64.exe main/main.go
endif

optimize:
	set CGO_ENABLED=0& set GOOS=windows& set GOARCH=amd64& go build -a -ldflags="-w -s" -o bin/main.exe main/main.go

docker-build:
	docker build -t kkt .

docker-run:
	docker run --name=kakao -p 8008:8008 kkt

docker-start:
	docker start kakao

docker-id:
	@docker ps -q --filter ancestor=kkt

docker-stop:
	@docker stop kakao

sync:
	aws s3 sync . s3://$(bucket)/chatbot

all: build run