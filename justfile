#!/usr/bin/env just --justfile

default: build

fmt:
    go fmt ./...

vet: fmt
    go vet ./...

build: vet
    go build -o aero ./cmd/aero

run:
    go run ./cmd/aero

test:
    go test ./internal...


update:
  go get -u
  go mod tidy -v
