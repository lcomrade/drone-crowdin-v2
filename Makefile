GO ?= go
GOFMT ?= gofmt

NAME = drone-crowdin-v2
MAIN_GO = ./cmd/*.go

LDFLAGS = -w -s -X "main.Version=$(shell git describe --tags --always | sed 's/-/+/' | sed 's/^v//')"


.PHONY: all fmt clean

all:
	mkdir -p ./dist/

	GOMODULE111=on CGO_ENABLED=0 $(GO) build -trimpath -ldflags="$(LDFLAGS)" -o ./dist/$(NAME) $(MAIN_GO)
	chmod +x ./dist/$(NAME)

fmt:
	@$(GOFMT) -w $(shell find ./ -type f -name '*.go')

clean:
	rm -rf ./dist/
