#
# Copyright (C) 2022, Xilinx Inc - All rights reserved
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

LIB_NAME := xilinx-container-runtime
LIB_VERSION := 0.0.1
GO_VERSION := $(shell go version | cut -c 14- | cut -d' ' -f1)
BUILD_TIME := $(shell date +"%Y-%m-%d")

# MODULE := xilinx.com/xilinx-container-runtime
MODULE := .
all: build

build: fmt check deps
	go build -o $(LIB_NAME) -ldflags "-s -w -X 'main.Version=$(LIB_VERSION)' -X 'main.GoVersion=$(GO_VERSION)' -X 'main.BuildTime=$(BUILD_TIME)'" $(MODULE)/src/cmd/...

# build deb package
debian: build
	mkdir -p $(LIB_NAME)_$(LIB_VERSION)_amd64/DEBIAN
	cp ./build/DEBIAN/control ./$(LIB_NAME)_$(LIB_VERSION)_amd64/DEBIAN
	mkdir -p $(LIB_NAME)_$(LIB_VERSION)_amd64/usr/bin
	cp ./$(LIB_NAME) ./$(LIB_NAME)_$(LIB_VERSION)_amd64/usr/bin
	mkdir -p $(LIB_NAME)_$(LIB_VERSION)_amd64/etc/$(LIB_NAME)
	cp ./src/configs/xilinx-container-runtime/config.toml ./$(LIB_NAME)_$(LIB_VERSION)_amd64/etc/$(LIB_NAME)
	dpkg-deb --build --root-owner-group $(LIB_NAME)_$(LIB_VERSION)_amd64
	rm -rf $(LIB_NAME)_$(LIB_VERSION)_amd64

# build rpm package
rpm: build
	mkdir -p ./build/rpm/SOURCES
	cp ./$(LIB_NAME) ./build/rpm/SOURCES
	cp ./src/configs/xilinx-container-runtime/config.toml ./build/rpm/SOURCES
	cd ./build/rpm && rpmbuild --define "_topdir `pwd`" -bb SPECS/xilinx-container-runtime.spec
	cp ./build/rpm/RPMS/x86_64/* ./

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
	cp ./src/configs/xilinx-container-runtime/config.toml /etc/xilinx-container-runtime/config.toml

clean:
	rm -rf ./xilinx-container-runtime*