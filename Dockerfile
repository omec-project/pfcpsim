# SPDX-License-Identifier: Apache-2.0
# Copyright 2022 Open Networking Foundation

# Stage pfcpsim-build: builds the pfcpsim docker image
FROM golang:alpine AS builder
WORKDIR /pfcpctl

COPY go.mod ./go.mod
COPY go.sum ./go.sum

RUN go mod download
# exploit local cache
VOLUME $(go env GOCACHE):/root/.cache/go-build

COPY . ./
RUN CGO_ENABLED=0 go build -o /bin/pfcpctl cmd/pfcpctl/main.go
RUN CGO_ENABLED=0 go build -o /bin/pfcpsim cmd/pfcpsim/main.go

# Stage pfcpsim: runtime image of pfcpsim, containing also pfcpctl
FROM alpine AS pfcpsim

COPY --from=builder /bin/pfcpctl /bin
COPY --from=builder /bin/pfcpsim /bin

RUN echo "export PATH=/bin:${PATH}" >> /root/.bashrc

ENTRYPOINT [ "pfcpsim" ]
