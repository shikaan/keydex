## build: Compile example to binary
.PHONY: build
build:
	@echo "> ($@) Running the build..."
	@mkdir -p .build
	@go build -v
	@echo "  ($@) Done!"

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

## info: generate package version - run by go:generate
.PHONY: info
info:
	@echo "> ($@) Generating build info..."
	@go run _scripts/info.go -version=${VERSION}
	@echo "  ($@) Done!"

## doc: Generate documentation - run by go:generate
.PHONY: docs
docs: info
	@echo "> ($@) Generating documentation..."
	@go run _scripts/docs.go
	@echo "  ($@) Done!"
