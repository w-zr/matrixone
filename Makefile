# Copyright 2021 - 2022 Matrix Origin
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#
# Examples
#
# By default, make builds the mo-service
#
# make
#
# To re-build MO -
#
#	make clean
#	make build
#
# To re-build MO in debug mode also with race detector enabled -
#
# make clean
# make debug
#
# To run static checks -
#
# make install-static-check-tools
# make static-check
#
# To construct a directory named vendor in the main module’s root directory that contains copies of all packages needed to support builds and tests of packages in the main module.
# make vendor
#

# where am I
ROOT_DIR = $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BIN_NAME := mo-service
UNAME_S := $(shell uname -s)
GOPATH := $(shell go env GOPATH)
GO_VERSION=$(shell go version)
BRANCH_NAME=$(shell git rev-parse --abbrev-ref HEAD)
LAST_COMMIT_ID=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date +%s)
MO_VERSION=$(shell git describe --always --contains $(shell git rev-parse HEAD))
GO_MODULE=$(shell go list -m)
MUSL_DIR=$(ROOT_DIR)/musl
MUSL_CC=$(MUSL_DIR)/bin/musl-gcc
MUSL_VERSION:=1.2.5

# cross compilation has been disabled for now
#ifneq ($(GOARCH)$(TARGET_ARCH)$(GOOS)$(TARGET_OS),)
#$(error cross compilation has been disabled)
#endif

###############################################################################
# default target
###############################################################################

all: build

###############################################################################
# build vendor directory
###############################################################################

.PHONY: vendor-build
vendor-build:
	$(info [go mod vendor])
	@go mod vendor

###############################################################################
# code generation
###############################################################################

.PHONY: config
config:
	$(info [Create build config])
	@go mod tidy

.PHONY: generate-pb
generate-pb:
	$(ROOT_DIR)/proto/gen.sh

# Generate protobuf files
.PHONY: pb
pb: vendor-build generate-pb fmt
	$(info all protos are generated)

###############################################################################
# build mo-service
###############################################################################

RACE_OPT :=
DEBUG_OPT :=
CGO_DEBUG_OPT :=
CGO_OPTS=CGO_CFLAGS="-I$(ROOT_DIR)/cgo " CGO_LDFLAGS="-L$(ROOT_DIR)/cgo -lm -lmo"
GOLDFLAGS=-ldflags="-X '$(GO_MODULE)/pkg/version.GoVersion=$(GO_VERSION)' -X '$(GO_MODULE)/pkg/version.BranchName=$(BRANCH_NAME)' -X '$(GO_MODULE)/pkg/version.CommitID=$(LAST_COMMIT_ID)' -X '$(GO_MODULE)/pkg/version.BuildTime=$(BUILD_TIME)' -X '$(GO_MODULE)/pkg/version.Version=$(MO_VERSION)'"
TAGS :=

# build mo-service binary
.PHONY: build
build: config cgo
	$(info [Build binary])
	$(CGO_OPTS) go build $(TAGS) $(RACE_OPT) $(GOLDFLAGS) $(DEBUG_OPT) -o $(BIN_NAME) ./cmd/mo-service

# build mo-tool
.PHONY: mo-tool
mo-tool: config cgo
	$(info [Build mo-tool tool])
	$(CGO_OPTS) go build -o mo-tool ./cmd/mo-tool

# build mo-service binary for debugging with go's race detector enabled
# produced executable is 10x slower and consumes much more memory
.PHONY: debug
debug: override BUILD_NAME := debug-binary
debug: override RACE_OPT := -race
debug: override DEBUG_OPT := -gcflags=all="-N -l"
debug: override CGO_DEBUG_OPT := debug
debug: build

###############################################################################
# clean
###############################################################################

.PHONY: clean
clean:
	$(info [Clean up])
	$(info Clean go test cache)
	@go clean -testcache
	rm -f $(BIN_NAME)
	rm -rf $(ROOT_DIR)/vendor
	rm -rf $(MUSL_DIR)
	rm -rf /tmp/musl-$(MUSL_VERSION).tar.gz
	$(MAKE) -C cgo clean
