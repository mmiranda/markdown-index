.PHONY: test
test:
	gotest -v ./...

build:
	go build
