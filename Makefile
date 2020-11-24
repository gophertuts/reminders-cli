.PHONY: client
.PHONY: server

all: fmt lint vet client server

fmt:
	@echo "Formatting the source code"
	go fmt ./...

lint:
	@echo "Linting the source code"
	# go get -u golang.org/x/lint/golint
	golint ./...

vet:
	@echo "Checking for code issues"
	go vet ./...

client:
	@echo "Removing the client binary"
	rm -f bin/client
	@echo "Building the client binary"
	go build -o bin/client cmd/client/main.go

server:
	@echo "Removing the server binary"
	rm -f bin/server
	@echo "Building the server binary"
	go build -o bin/server cmd/server/main.go
