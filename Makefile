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

BIN_ROOT_DIR=examples
BIN_ID=${MODULE_ROOT}/${BIN_ROOT_DIR}
BIN_SRCS=\
	${BIN_ROOT_DIR}/mdnssearch \
	${BIN_ROOT_DIR}/mdnsserver
BINS=\
	${BIN_ID}/mdnssearch \
	${BIN_ID}/mdnsserver

.PHONY: format vet lint clean

all: test

format:
	gofmt -s -w ${PKG_SRC_DIR} ${TEST_PKG_DIR} ${BIN_ROOT_DIR}

vet: format
	go vet ${PKG_ID} ${TEST_PKG_ID} ${BINS}

lint: format
	golangci-lint run ${PKG_SRC_DIR}/... ${TEST_PKG_DIR}/...

test: lint
	go test -v -p 1 -timeout 10m -cover -coverpkg=${PKG}/... -coverprofile=${PKG_COVER}.out ${PKG}/... ${TEST_PKG}/...
	go tool cover -html=${PKG_COVER}.out -o ${PKG_COVER}.html

install: test
	go install ${BINS}

clean:
	go clean -i ${PKG}  ${TEST_PKG} ${BINS}
