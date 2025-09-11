# SPDX-License-Identifier: Apache-2.0
# Copyright 2022-present Open Networking Foundation
# Copyright 2024-present Intel Corporation

# Stage pfcpsim-build: builds the pfcpsim docker image
FROM golang:1.25.1-bookworm AS builder
WORKDIR /pfcpctl

COPY . .

RUN CGO_ENABLED=0 go build -o ./pfcpctl cmd/pfcpctl/main.go && \
    CGO_ENABLED=0 go build -o ./pfcpsim cmd/pfcpsim/main.go

# Stage pfcpsim: runtime image of pfcpsim, containing also pfcpctl
FROM alpine:3.22 AS pfcpsim

RUN apk update && apk add --no-cache -U tcpdump

COPY --from=builder /pfcpctl/pfcp* /usr/local/bin

ENTRYPOINT [ "pfcpsim" ]
