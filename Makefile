# Copyright (C) 2022 The go-mdns Authors All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SHELL := bash

PATH := $(GOBIN):$(PATH)
GOBIN := $(shell go env GOPATH)/bin

MODULE_ROOT=github.com/cybergarage/go-mdns

PKG_NAME=mdns
PKG_VER=$(shell git describe --abbrev=0 --tags)
PKG_COVER=${PKG_NAME}-cover

PKG_ID=${MODULE_ROOT}/${PKG_NAME}
PKG_SRC_DIR=${PKG_NAME}
PKG=${MODULE_ROOT}/${PKG_SRC_DIR}

TEST_PKG_NAME=${PKG_NAME}test
TEST_PKG_ID=${MODULE_ROOT}/${TEST_PKG_NAME}
TEST_PKG_DIR=${TEST_PKG_NAME}
TEST_PKG=${MODULE_ROOT}/${TEST_PKG_DIR}

BIN_ROOT_DIR=cmd
BIN_ID=${MODULE_ROOT}/${BIN_ROOT_DIR}
BIN_LOOKUP=${PKG_NAME}lookup
BIN_SERVER=${PKG_NAME}d
BIN_SRCS=\
	${BIN_ROOT_DIR}/${BIN_LOOKUP} \
	${BIN_ROOT_DIR}/${BIN_SERVER}
BINS=\
	${BIN_ID}/${BIN_LOOKUP} \
	${BIN_ID}/${BIN_SERVER}

DOCS_ROOT_DIR=doc

.PHONY: format vet lint clean
.IGNORE: lint

all: test

# The version is generated from the latest tag. It is kept as it is when there
# is no tag yet, and when the generated version is not newer than the one which
# version.go holds, so that a release which is prepared by hand is not lost.
version:
	@cur_version=$$(sed -n 's/.*Version = "\(.*\)".*/\1/p' ${PKG_SRC_DIR}/version.go); \
	new_version=$$(cd ${PKG_SRC_DIR} && ./version.gen 2>/dev/null | sed -n 's/.*Version = "\(.*\)".*/\1/p'); \
	if [ -z "$$new_version" ]; then \
		echo "version: no tag is found: $$cur_version is kept"; \
	elif [ "$$cur_version" = "$$new_version" ]; then \
		echo "version: $$cur_version is up to date"; \
	elif [ "$$(printf '%s\n%s\n' "$$cur_version" "$$new_version" | sort -V | tail -n 1)" != "$$new_version" ]; then \
		echo "version: $$cur_version is kept ($$new_version is not newer)"; \
	else \
		(cd ${PKG_SRC_DIR} && ./version.gen > version.go) && \
		echo "version: $$cur_version -> $$new_version"; \
		git commit ${PKG_SRC_DIR}/version.go -m "Update version" || true; \
	fi

format:version
	gofmt -s -w ${PKG_SRC_DIR} ${TEST_PKG_DIR} ${BIN_ROOT_DIR}

vet: format
	go vet ${PKG_ID} ${TEST_PKG_ID} ${BINS}

lint: format
	golangci-lint run ${PKG_SRC_DIR}/... ${TEST_PKG_DIR}/...

test: lint
	go test -v -race -p 1 -timeout 10m -cover -coverpkg=${PKG}/... -coverprofile=${PKG_COVER}.out ${PKG}/... ${TEST_PKG}/...
	go tool cover -html=${PKG_COVER}.out -o ${PKG_COVER}.html

install:
	go install ${BINS}
	${GOBIN}/${BIN_LOOKUP} doc > ${DOCS_ROOT_DIR}/${BIN_LOOKUP}.md
	@git diff --quiet -- ${DOCS_ROOT_DIR}/${BIN_LOOKUP}.md || \
		git commit ${DOCS_ROOT_DIR}/${BIN_LOOKUP}.md -m "docs: update ${BIN_LOOKUP} command reference"
# ${BIN_SERVER} is under development, and it has no doc command yet, so
# ${DOCS_ROOT_DIR}/${BIN_SERVER}.md is maintained by hand.

clean:
	go clean -i ${PKG} ${TEST_PKG} ${BINS}
