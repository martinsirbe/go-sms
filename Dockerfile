FROM golang:1.14-alpine3.11 AS builder

RUN apk update && apk add git ca-certificates

WORKDIR /go/src/github.com/martinsirbe/go-sms
COPY . .

ENV GO111MODULE=on
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o go-sms ./cmd/go-sms

FROM alpine:3.11

COPY --from=builder /go/src/github.com/martinsirbe/go-sms/go-sms /bin/go-sms
