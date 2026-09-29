BINARY_NAME=run-nothing

.PHONY: all build run interactive test clean release

all: build

build:
	go build -ldflags="-s -w" -o $(BINARY_NAME) .

run:
	go run .

interactive:
	go run . --interactive

test:
	go test -v ./...

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)-*

release:
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BINARY_NAME)-macos-arm64 .
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BINARY_NAME)-macos-x86_64 .
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BINARY_NAME)-linux-x86_64 .
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $(BINARY_NAME)-linux-arm64 .
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BINARY_NAME)-windows-x86_64.exe .
