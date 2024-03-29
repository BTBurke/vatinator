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
NEXTVERSION = $$(( $$(git tag --sort=-v:refname | head -n 1 | cut -d. -f1 | sed 's/v//') + 1 ))

build: bundled/assets.go
> mkdir -p bin
> go build -o bin/vat $$(find cmd/vat -name '*.go')

test: bundled/assets.go
> go test ./...

assets:
> go get -u github.com/go-bindata/go-bindata/...
> go-bindata -pkg bundled -o bundled/assets.go assets/

bundled/assets.go: $(shell find ./assets -type f)
> go get -u github.com/go-bindata/go-bindata/...
> go-bindata -pkg bundled -o bundled/assets.go assets/

encrypt:
> go build -o enc build/encrypt.go
> ./enc vatinator-f91ccb107c2c.json
> rm ./enc

tag:
> git tag -a v${NEXTVERSION}.0.0 -m "v${NEXTVERSION}.0.0"

test-release:
> goreleaser build --snapshot --rm-dist

deploy:
> go build -o server cmd/server/main.go
> mv server /usr/local/bin
> systemctl restart vat.service
> systemctl status vat.service

.PHONY: build test assets encrypt tag test-release deploy
