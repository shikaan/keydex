## build: Compile example to binary
.PHONY: build
build: docs
	@echo "Running go generate to prepare build info..."
	@go generate
	@echo "Running the build..."
	@go build -v
	@echo "Done!"

## info: generate package version - run by go:generate
.PHONY: info
info:
	@echo "> Generating build info..."
	@go run _scripts/info.go

## doc: Generate documentation
.PHONY: docs
docs:
	@echo "> Generating documentation..."
	@go run _scripts/docs.go

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
