PROJECT_NAME := go-sms
GOLANGCI_LINT_VER := v1.24

.PHONY: go-gen
go-gen:
	@go generate ./...

.PHONY: tests
tests:
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

.PHONY: go-lint
go-lint:
	docker run \
		--rm \
		--volume "$$PWD":/src/github.com/martinsirbe/$(PROJECT_NAME) \
		--workdir /src/github.com/martinsirbe/$(PROJECT_NAME) \
		golangci/golangci-lint:$(GOLANGCI_LINT_VER) \
		/bin/bash -c "golangci-lint run -v --config=.golangci.yml"
