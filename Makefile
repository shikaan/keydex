REVISION=git log -n1 --pretty=%h
VERSION=git describe --abbrev=0 --tags 2> /dev/null || echo dev
NAME=echo "kpcli"

## info: generate package version - run by go:generate
.PHONY: info
info:
	@sed "s/_REVISION_/`$(REVISION)`/; s/_VERSION_/`$(VERSION)`/; s/_NAME_/`$(NAME)`/" ./pkg/info/info.tmpl > ./pkg/info/info.go

## run: Executes the entrypoint
.PHONY: run
run:
	go generate
	go run . edit test.kdbx /another_test

## build: Compile example to binary
.PHONY: build
build:
	go generate -x
	go build -v

## clean: Clean compilation artifacts
.PHONY: clean
clean:
	@echo "Nothing to clean yet!"

## test: Run tests
.PHONY: test
test:
	@echo "No tests yet!"

## help: Show help and exit
.PHONY: help
help: Makefile
	@echo
	@echo "  Choose a command:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo