
LIB_NAME := xilinx-container-runtime
LIB_VERSION := 0.0.1
GO_VERSION := $(shell go version | cut -c 14- | cut -d' ' -f1)
BUILD_TIME := $(shell date +"%Y-%m-%d")
PKG_REV := 1

# MODULE := xilinx.com/xilinx-container-runtime
MODULE := .
all: build

build: fmt check deps
	go build -o xilinx-container-runtime -ldflags "-s -w -X 'main.Version=$(LIB_VERSION)' -X 'main.GoVersion=$(GO_VERSION)' -X 'main.BuildTime=$(BUILD_TIME)'" $(MODULE)/cmd/...

# Define the check targets for the Golang codebase
.PHONY: check fmt assert-fmt ineffassign lint vet
check: assert-fmt vet
fmt:
	go list -f '{{.Dir}}' $(MODULE)/... \
		| xargs gofmt -s -l -w

assert-fmt:
	go list -f '{{.Dir}}' $(MODULE)/... \
		| xargs gofmt -s -l > fmt.out
	@if [ -s fmt.out ]; then \
		echo "\nERROR: The following files are not formatted:\n"; \
		cat fmt.out; \
		rm fmt.out; \
		exit 1; \
	else \
		rm fmt.out; \
	fi

vet:
	go vet $(MODULE)/...

deps:
	go mod download

install:
	cp ./xilinx-container-runtime /usr/bin/xilinx-container-runtime
	@if [ -d "/etc/xilinx-container-runtime" ]; then \
        echo "Dir /etc/xilinx-container-runtime existed"; \
	else \
		mkdir /etc/xilinx-container-runtime; \
    fi
	cp ./configs/xilinx-container-runtime/config.toml /etc/xilinx-container-runtime/config.toml
	# cp ./configs/docker/daemon.json /etc/docker/daemon.json

clean:
	rm ./xilinx-container-runtime