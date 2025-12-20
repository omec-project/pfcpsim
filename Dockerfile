# SPDX-License-Identifier: Apache-2.0
# Copyright 2022-present Open Networking Foundation
# Copyright 2024-present Intel Corporation

# Stage pfcpsim-build: builds the pfcpsim docker image
FROM golang:1.25.5-bookworm@sha256:09f53deea14d4019922334afe6258b7b776afc1d57952be2012f2c8c4076db05 AS builder
WORKDIR /pfcpctl

COPY . .

RUN CGO_ENABLED=0 go build -o ./pfcpctl cmd/pfcpctl/main.go && \
    CGO_ENABLED=0 go build -o ./pfcpsim cmd/pfcpsim/main.go

# Stage pfcpsim: runtime image of pfcpsim, containing also pfcpctl
FROM alpine:3.23@sha256:865b95f46d98cf867a156fe4a135ad3fe50d2056aa3f25ed31662dff6da4eb62 AS pfcpsim

RUN apk update && apk add --no-cache -U tcpdump

COPY --from=builder /pfcpctl/pfcp* /usr/local/bin

ENTRYPOINT [ "pfcpsim" ]
