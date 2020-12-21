SHELL := /usr/bin/bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later)
endif
.RECIPEPREFIX = >

build:
> mkdir -p bin
> go build -o bin/vat $$(find cmd/vat -name '*.go')

test:
> go test ./...

assets: encrypt
> go get -u github.com/go-bindata/go-bindata/...
> go-bindata -pkg bundled -o bundled/assets.go assets/

encrypt:
> go build -o enc build/encrypt.go
> ./enc vatinator-f91ccb107c2c.json
> rm ./enc

.PHONY: build test assets encrypt
