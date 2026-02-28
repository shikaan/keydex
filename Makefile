## build: Make a release binary
.PHONY: build
build:
	@echo "> ($@) Creating a build..."
	@mkdir -p .build
	@go build -v
	@echo "  ($@) Done!"

## build-debug: Make a debug binary
.PHONY: build-debug
build-debug:
	@echo "> ($@) Creating a debug build..."
	@mkdir -p .build
	@go build -gcflags=all="-N -l" -v
	@echo "  ($@) Done!"

## test: Run tests
.PHONY: test
test: info build
	@go test ./...

## coverage: Run coverage
.PHONY: coverage
coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

## help: Show help and exit
.PHONY: help
help: Makefile
	@echo
	@echo "  Choose a command:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## info: generate package version - run by go generate
.PHONY: info
info:
	@echo "> ($@) Generating build info..."
	@go run info.go -version=${VERSION}
	@echo "  ($@) Done!"

## docs: Generate documentation - run by go generate
.PHONY: docs
docs: info
	@echo "> ($@) Generating documentation..."
	@go run docs.go
	@echo "  ($@) Done!"

## clean: Remove build artifacts
.PHONY: clean
clean:
	@rm -rf keydex
	@rm -rf .build/**
