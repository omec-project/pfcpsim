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

DOCKER_TARGET           ?= pfcpsim-client

run-pfcpsim-client:
	go run cmd/pfcpsim-client/main.go -v -i en0 -c 2

build-pfcpsim-client:
	DOCKER_BUILDKIT=$(DOCKER_BUILDKIT) docker build -f cmd/pfcpsim-client/Dockerfile . \
	--target ${DOCKER_TARGET} \
	--tag ${DOCKER_REGISTRY}${DOCKER_REPOSITORY}${DOCKER_TARGET}:${DOCKER_TAG} \
