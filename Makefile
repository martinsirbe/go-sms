PROJECT_NAME := go-sms
GOLANGCI_LINT_VER := v1.17.1

.PHONY: go-gen
go-gen:
	@go generate ./...

.PHONY: test
test:
	@go test -v -cover -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: build
build:
	@go build -o bin/go-sms cmd/go-sms/main.go

.PHONY: lint
lint:
	@docker run --rm -w /src/github.com/martinsirbe/$(PROJECT_NAME) \
	    -v "$$PWD":/src/github.com/martinsirbe/$(PROJECT_NAME) \
	     golangci/golangci-lint:$(GOLANGCI_LINT_VER) golangci-lint run -v
