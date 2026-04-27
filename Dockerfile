# SPDX-License-Identifier: Apache-2.0
# Copyright 2022-present Open Networking Foundation
# Copyright 2024-present Intel Corporation

# Stage pfcpsim-build: builds the pfcpsim docker image
FROM golang:1.26.2-bookworm@sha256:47ce5636e9936b2c5cbf708925578ef386b4f8872aec74a67bd13a627d242b19 AS builder
WORKDIR /pfcpctl

COPY . .

RUN CGO_ENABLED=0 go build -o ./pfcpctl cmd/pfcpctl/main.go && \
    CGO_ENABLED=0 go build -o ./pfcpsim cmd/pfcpsim/main.go

# Stage pfcpsim: runtime image of pfcpsim, containing also pfcpctl
FROM alpine:3.23@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11 AS pfcpsim

# Build arguments for dynamic labels
ARG VERSION=dev
ARG VCS_URL=unknown
ARG VCS_REF=unknown
ARG BUILD_DATE=unknown

LABEL org.opencontainers.image.source="${VCS_URL}" \
    org.opencontainers.image.version="${VERSION}" \
    org.opencontainers.image.created="${BUILD_DATE}" \
    org.opencontainers.image.revision="${VCS_REF}" \
    org.opencontainers.image.url="${VCS_URL}" \
    org.opencontainers.image.title="pfcpsim" \
    org.opencontainers.image.description="Aether 5G Core PFCPSIM Network Function" \
    org.opencontainers.image.authors="Aether SD-Core <dev@lists.aetherproject.org>" \
    org.opencontainers.image.vendor="Aether Project" \
    org.opencontainers.image.licenses="Apache-2.0" \
    org.opencontainers.image.documentation="https://docs.sd-core.aetherproject.org/"

RUN apk add --no-cache tcpdump

COPY --from=builder /pfcpctl/pfcp* /usr/local/bin

ENTRYPOINT [ "pfcpsim" ]
