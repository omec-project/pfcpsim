# SPDX-License-Identifier: Apache-2.0
# Copyright 2022-present Open Networking Foundation
# Copyright 2024-present Intel Corporation

# Stage pfcpsim-build: builds the pfcpsim docker image
FROM golang:1.26-bookworm@sha256:eae3cdfa040d0786510a5959d36a836978724d03b34a166ba2e0e198baac9196 AS builder
WORKDIR /pfcpctl

COPY . .

RUN CGO_ENABLED=0 go build -o ./pfcpctl cmd/pfcpctl/main.go && \
    CGO_ENABLED=0 go build -o ./pfcpsim cmd/pfcpsim/main.go

# Stage pfcpsim: runtime image of pfcpsim, containing also pfcpctl
FROM alpine:3.23@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659 AS pfcpsim

RUN apk update && apk add --no-cache -U tcpdump

COPY --from=builder /pfcpctl/pfcp* /usr/local/bin

ENTRYPOINT [ "pfcpsim" ]
