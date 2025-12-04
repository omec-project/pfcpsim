# SPDX-License-Identifier: Apache-2.0
# Copyright 2022-present Open Networking Foundation
# Copyright 2024-present Intel Corporation

# Stage pfcpsim-build: builds the pfcpsim docker image
FROM golang:1.25.4-bookworm@sha256:e17419604b6d1f9bc245694425f0ec9b1b53685c80850900a376fb10cb0f70cb AS builder
WORKDIR /pfcpctl

COPY . .

RUN CGO_ENABLED=0 go build -o ./pfcpctl cmd/pfcpctl/main.go && \
    CGO_ENABLED=0 go build -o ./pfcpsim cmd/pfcpsim/main.go

# Stage pfcpsim: runtime image of pfcpsim, containing also pfcpctl
FROM alpine:3.23@sha256:51183f2cfa6320055da30872f211093f9ff1d3cf06f39a0bdb212314c5dc7375 AS pfcpsim

RUN apk update && apk add --no-cache -U tcpdump

COPY --from=builder /pfcpctl/pfcp* /usr/local/bin

ENTRYPOINT [ "pfcpsim" ]
