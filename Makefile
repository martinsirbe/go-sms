.PHONY: go-gen
go-gen:
	@go generate ./...

.PHONY: test
test:
	@go test -v --cover ./...

.PHONY: build
build:
	@go build -o bin/go-sms cmd/go-sms/main.go

.PHONY: lint
lint:
	golangci-lint run -v
