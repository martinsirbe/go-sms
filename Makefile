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

.PHONY: docker-build
docker-build:
	@docker build -t go-sms -f Dockerfile .

.PHONY: run
run:
	docker run --rm -v $$(pwd):/root/home go-sms:latest /bin/go-sms --config-path=/root/home/config.yaml

.PHONY: lint
lint:
	@docker run --rm -w /src/github.com/martinsirbe/$(PROJECT_NAME) \
	    -v "$$PWD":/src/github.com/martinsirbe/$(PROJECT_NAME) \
	     golangci/golangci-lint:$(GOLANGCI_LINT_VER) golangci-lint run -v
