# SPDX-License-Identifier: Apache-2.0
# Copyright 2022-present Open Networking Foundation

PROJECT_NAME             := pfcpsim
VERSION                  ?= $(shell cat ./VERSION)

## Docker related
DOCKER_REGISTRY          ?=
DOCKER_REPOSITORY        ?=
DOCKER_TAG               ?= ${VERSION}
DOCKER_IMAGENAME         := ${DOCKER_REGISTRY}${DOCKER_REPOSITORY}${PROJECT_NAME}:${DOCKER_TAG}
DOCKER_BUILDKIT          ?= 1

DOCKER_TARGET           ?= pfcpsim

build-pfcpsim:
	DOCKER_BUILDKIT=$(DOCKER_BUILDKIT) docker build -f Dockerfile . \
    	--target $(DOCKER_TARGET) \
    	--cache-from ${DOCKER_REGISTRY}${DOCKER_REPOSITORY}$(DOCKER_TARGET):${DOCKER_TAG} \
    	--tag ${DOCKER_REGISTRY}${DOCKER_REPOSITORY}$(DOCKER_TARGET):${DOCKER_TAG}

golint:
	@docker run --rm -v $(CURDIR):/app -w /app/pkg/pfcpsim golangci/golangci-lint:latest golangci-lint run -v --config /app/.golangci.yml

.coverage:
	rm -rf $(CURDIR)/.coverage
	mkdir -p $(CURDIR)/.coverage

test: .coverage
	go test	-race -coverprofile=.coverage/coverage-unit.txt -covermode=atomic -v ./...

build-proto:
	@echo "Compiling proto files..."
	docker run --rm -v $(CURDIR)/api:/source diebietse/go-gw-protoc:latest \
    -I./ \
    -I/usr/local/include \
    --go_out=paths=source_relative,plugins=grpc:./ \
    pfcpsim.proto
