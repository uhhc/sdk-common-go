OS = Linux
VERSION = 0.0.1

CURDIR = $(shell pwd)
SOURCEDIR = $(CURDIR)
COVER = $($3)

ECHO = echo
RM = rm -rf
MKDIR = mkdir

.PHONY: test grpc

default: test lint vet

test: fmt lint
	go test -cover=true $(PACKAGES)

race:
	go test -cover=true -race $(PACKAGES)

# http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources
fmt:
	go fmt ./...

# https://godoc.org/golang.org/x/tools/cmd/goimports
# imports:
# 	goimports ./...

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint:
	golint ./...

# http://godoc.org/code.google.com/p/go.tools/cmd/vet
# go get code.google.com/p/go.tools/cmd/vet
vet:
	go vet ./...

all: test

PACKAGES = $(shell go list ./... | grep -v './vendor/\|./tests\|./mock')
BUILD_PATH = $(shell if [ "$(CI_DEST_DIR)" != "" ]; then echo "$(CI_DEST_DIR)" ; else echo "$(PWD)"; fi)

cover: collect-cover-data test-cover-html open-cover-html

collect-cover-data:
	echo "mode: count" > coverage-all.out
	@$(foreach pkg,$(PACKAGES),\
		go test -v -coverprofile=coverage.out -covermode=count $(pkg);\
		if [ -f coverage.out ]; then\
			tail -n +2 coverage.out >> coverage-all.out;\
		fi;)

test-cover-html:
	go tool cover -html=coverage-all.out -o coverage.html

test-cover-func:
	go tool cover -func=coverage-all.out

open-cover-html:
	open coverage.html

run: fmt lint
	go run main.go

clean:
	rm -f *.out *.html

grpc:
	# Format: protoc -I <proto_import_path1> -I <proto_import_path2> --go_out=<plugins={plugin1+plugin2+...}>:<pb_output_path> <proto_file_path>
# 	protoc -I pkg/grpc/proto --gogofaster_out=plugins=grpc:pkg/grpc/pb pkg/grpc/proto/*.proto
	protoc -I grpc/proto --go_out=plugins=grpc:grpc/pb grpc/proto/*.proto

compile: test build

help:
	@$(ECHO) "Targets:"
	@$(ECHO) "all				- test"
	@$(ECHO) "setup				- install necessary libraries"
	@$(ECHO) "test				- run all unit tests"
	@$(ECHO) "cover [package]	- generates and opens unit test coverage report for a package"
	@$(ECHO) "race				- run all unit tests in race condition"
	@$(ECHO) "add				- runs govendor add +external command"
	@$(ECHO) "build				- build and exports using CI_DEST_DIR"
	@$(ECHO) "build-local       - build and exports using CI_DEST_DIR locally"
	@$(ECHO) "run				- run the main.go"
	@$(ECHO) "clean				- remove test reports and compiled package from this folder"
	@$(ECHO) "compile			- test and build - one command for CI"
	@$(ECHO) "grpc		        - init or rebuild protocol buffers files"