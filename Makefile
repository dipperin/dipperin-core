# Dipperin

# project name
PROJECTNAME=$(shell basename "$(PWD)")

# project path
ROOT=$(shell pwd)

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## test: run dipperin all unit tests
test:
	@echo "run Dipperin unit test"
	@./cs.sh

## build: build dipperin-core all cmd to $GOPATH/bin
build:
	@echo "build Dipperin-core cmd"
	@./cs.sh install

## tidy: update vendor pkg, also support http_proxy, you need add port=you_http_proxy_port
tidy:
	@echo "start update vendor"
	@./cs.sh tidy $(port)

## cover: get test coverage
cover:
	@./cs.sh cover $(show)

## travis-test
travis-test:
	@go test ./...

## cross compiling
compile:
	@./cs.sh compile
