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
PKG_ID=${MODULE_ROOT}/${PKG_NAME}
PKG_SRC_DIR=${PKG_NAME}
PKG_SRCS=\
	${PKG_SRC_DIR} \
	${PKG_SRC_DIR}/encoding \
	${PKG_SRC_DIR}/protocol \
	${PKG_SRC_DIR}/transport
PKGS=\
	${PKG_ID} \
	${PKG_ID}/encoding \
	${PKG_ID}/protocol \
	${PKG_ID}/transport

TEST_PKG_NAME=${PKG_NAME}test
TEST_PKG_ID=${MODULE_ROOT}/${TEST_PKG_NAME}
TEST_PKG_DIR=${TEST_PKG_NAME}
TEST_PKG_SRCS=\
	${TEST_PKG_DIR}
TEST_PKGS=\
	${TEST_PKG_ID}

BIN_ROOT=examples
BIN_ID=${MODULE_ROOT}/${BIN_ROOT}
BIN_SRCS=\
	${BIN_ROOT}/mdnssearch \
	${BIN_ROOT}/mdnsserver
BINS=\
	${BIN_ID}/mdnssearch \
	${BIN_ID}/mdnsserver

.PHONY: format vet lint clean

all: test

format:
	gofmt -s -w ${PKG_SRC_DIR} ${TEST_PKG_DIR}

vet: format
	go vet ${PKG_ID} ${TEST_PKG_ID}

lint: format
	golangci-lint run ${PKG_SRCS} ${TEST_PKG_SRCS} ${BIN_SRCS}

test: lint
	go test -v -cover -timeout 60s ${PKGS} ${TEST_PKGS}

install: test
	go install ${BINS}

clean:
	go clean -i ${PKGS}  ${TEST_PKGS} ${BINS}
