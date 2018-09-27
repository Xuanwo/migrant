SHELL := /bin/bash

PROJECT_NAME:=migrant

.PHONY: all check format　vet lint build install uninstall release clean test coverage

VERSION=$(shell cat ./constants/version.go | grep "Version\ =" | sed -e s/^.*\ //g | sed -e s/\"//g)
DIRS_TO_CHECK=$(shell ls -d */ | grep -vE "vendor|test")
PKGS_TO_CHECK=$(shell go list ./... | grep -v "/vendor/")

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  check      to format, vet and lint "
	@echo "  build      to create bin directory and build migrant"
	@echo "  install    to install migrant to /usr/local/bin/migrant"
	@echo "  uninstall  to uninstall migrant"
	@echo "  release    to release migrant"
	@echo "  clean      to clean build and test files"
	@echo "  test       to run test"
	@echo "  coverage   to test with coverage"

check: format vet lint

format:
	@echo "vgo fmt, skipping vendor packages"
	@for pkg in ${PKGS_TO_CHECK}; do vgo fmt $${pkg}; done;
	@echo "ok"

vet:
	@echo "vgo tool vet, skipping vendor packages"
	@vgo tool vet -all ${DIRS_TO_CHECK}
	@echo "ok"

lint:
	@echo "golint, skipping vendor packages"
	@lint=$$(for pkg in ${PKGS_TO_CHECK}; do golint $${pkg}; done); \
	 lint=$$(echo "$${lint}"); \
	 if [[ -n $${lint} ]]; then echo "$${lint}"; exit 1; fi
	@echo "ok"

build: check
	@echo "build migrant"
	@mkdir -p ./bin
	@vgo build -o ./bin/migrant .
	@echo "ok"

install: build
	@echo "install migrant to GOPATH"
	@cp ./bin/migrant ${GOPATH}/bin/migrant
	@echo "ok"

uninstall:
	@echo "delete /usr/local/bin/migrant"
	@rm -f /usr/local/bin/migrant
	@echo "ok"

release:
	@echo "release migrant"
	@rm ./release/*
	@mkdir -p ./release

	@echo "build for linux"
	@GOOS=linux GOARCH=amd64 vgo build -o ./bin/linux/migrant_v${VERSION}_linux_amd64 .
	@tar -C ./bin/linux/ -czf ./release/migrant_v${VERSION}_linux_amd64.tar.gz migrant_v${VERSION}_linux_amd64

	@echo "build for macOS"
	@GOOS=darwin GOARCH=amd64 vgo build -o ./bin/macos/migrant_v${VERSION}_macos_amd64 .
	@tar -C ./bin/macos/ -czf ./release/migrant_v${VERSION}_macos_amd64.tar.gz migrant_v${VERSION}_macos_amd64

	@echo "build for windows"
	@GOOS=windows GOARCH=amd64 vgo build -o ./bin/windows/migrant_v${VERSION}_windows_amd64.exe .
	@tar -C ./bin/windows/ -czf ./release/migrant_v${VERSION}_windows_amd64.tar.gz migrant_v${VERSION}_windows_amd64.exe

	@echo "ok"

clean:
	@rm -rf ./bin
	@rm -rf ./release
	@rm -rf ./coverage

test:
	@echo "run test"
	@vgo test -v ${PKGS_TO_CHECK}
	@echo "ok"

coverage:
	@echo "run test with coverage"
	@for pkg in ${PKGS_TO_CHECK}; do \
		output="coverage$${pkg#github.com/migrant/migrant}"; \
		mkdir -p $${output}; \
		vgo test -v -cover -coverprofile="$${output}/profile.out" $${pkg}; \
		if [[ -e "$${output}/profile.out" ]]; then \
			vgo tool cover -html="$${output}/profile.out" -o "$${output}/profile.html"; \
		fi; \
	done
	@echo "ok"

.PHONY: start
start: build
	@echo "Starting server..."
	@./bin/migrant -c ${CONFIG_FILE}
	@echo "Done"
