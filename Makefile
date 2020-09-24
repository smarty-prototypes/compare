#!/usr/bin/make -f

test: fmt
	go test -count=1 -short -cover $(ARGS) ./...

fmt:
	@go version
	go mod tidy
	go fmt ./...
