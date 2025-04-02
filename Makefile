# SPDX-License-Identifier: Apache-2.0
# Copyright 2022-present Open Networking Foundation

PROJECT_NAME             := sdcore
DOCKER_VERSION           ?= $(shell cat ./VERSION)

## Docker related
DOCKER_REGISTRY          ?=
DOCKER_REPOSITORY        ?=
DOCKER_TAG               ?= ${DOCKER_VERSION}
DOCKER_IMAGENAME         := ${DOCKER_REGISTRY}${DOCKER_REPOSITORY}${PROJECT_NAME}:${DOCKER_TAG}
DOCKER_BUILDKIT          ?= 1

DOCKER_TARGET            ?= pfcpsim

docker-build:
	DOCKER_BUILDKIT=$(DOCKER_BUILDKIT) docker build \
		--target $(DOCKER_TARGET) \
		--tag ${DOCKER_REGISTRY}${DOCKER_REPOSITORY}$(DOCKER_TARGET):${DOCKER_TAG} \
		.

docker-push:
	docker push ${DOCKER_REGISTRY}${DOCKER_REPOSITORY}$(DOCKER_TARGET):${DOCKER_TAG}

golint:
	@docker run --rm -v $(CURDIR):/app -w /app/pkg/pfcpsim golangci/golangci-lint:latest golangci-lint run -v --config /app/.golangci.yml

.coverage:
	rm -rf $(CURDIR)/.coverage
	mkdir -p $(CURDIR)/.coverage

# -run flag ensures that the fuzz test won't be run
# because the fuzz test needs a UPF to run
test: .coverage
	go test	-race -coverprofile=.coverage/coverage-unit.txt -covermode=atomic -run=^Test -v ./...

reuse-lint:
	docker run --rm -v $(CURDIR):/pfcpsim -w /pfcpsim omecproject/reuse-verify:latest reuse lint

docker-protobuf:
	DOCKER_BUILDKIT=$(DOCKER_BUILDKIT) docker build \
		--tag protobuf:latest \
		-f $(CURDIR)/api/Dockerfile \
		.

build-proto: docker-protobuf
	@echo "Compiling proto files..."
	docker run --rm -v $(CURDIR)/api:/source -w /source protobuf:latest \
		protoc -I./ \
		--go_out=paths=source_relative:./ \
		--go-grpc_out=paths=source_relative:./ \
		pfcpsim.proto
