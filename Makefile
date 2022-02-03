# SPDX-License-Identifier: Apache-2.0
# Copyright 2022-present Open Networking Foundation

PROJECT_NAME             := pfcpsim
VERSION                  ?= $(shell cat ./VERSION)

# tool containers
VOLTHA_TOOLS_VERSION ?= 2.3.1

## Docker related
DOCKER_REGISTRY          ?=
DOCKER_REPOSITORY        ?=
DOCKER_TAG               ?= ${VERSION}
DOCKER_IMAGENAME         := ${DOCKER_REGISTRY}${DOCKER_REPOSITORY}${PROJECT_NAME}:${DOCKER_TAG}
DOCKER_BUILDKIT          ?= 1

DOCKER_TARGET           ?= pfcpsim-client

build-pfcpsim-client:
	DOCKER_BUILDKIT=$(DOCKER_BUILDKIT) docker build -f Dockerfile . \
	--target ${DOCKER_TARGET} \
	--tag ${DOCKER_REGISTRY}${DOCKER_REPOSITORY}${DOCKER_TARGET}:${DOCKER_TAG}

golint:
	@docker run --rm -v $(CURDIR):/app -w /app/pkg/pfcpsim golangci/golangci-lint:latest golangci-lint run -v --config /app/.golangci.yml

.coverage:
	rm -rf $(CURDIR)/.coverage
	mkdir -p $(CURDIR)/.coverage

test: .coverage
	go test	-race -coverprofile=.coverage/coverage-unit.txt -covermode=atomic -v ./...

build-proto:
	@echo "Compiling proto files..."
	 docker run --rm -v ${PWD}/api:/source diebietse/go-gw-protoc -I. --go_out=plugins=grpc:. --grpc-gateway_out=logtostderr=true:. pfcpsim.proto

