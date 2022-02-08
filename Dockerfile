# SPDX-License-Identifier: Apache-2.0
# Copyright 2022 Open Networking Foundation

# Stage pfcpsim-build: builds the pfcpsim docker image
FROM golang:alpine AS builder
WORKDIR /pfcpsimctl

COPY go.mod ./go.mod
COPY go.sum ./go.sum

RUN go mod download
# exploit local cache
VOLUME $(go env GOCACHE):/root/.cache/go-build

COPY . ./
RUN CGO_ENABLED=0 go build -o /bin/pfcpsimctl cmd/pfcpsimctl/main.go
RUN CGO_ENABLED=0 go build -o /bin/pfcpsim cmd/pfcpsim/main.go

# Stage pfcpsimctl: runtime image of pfcpsim, containing also pfcpsimctl
FROM golang:alpine AS pfcpsim

COPY --from=builder /bin/pfcpsimctl /bin
COPY --from=builder /bin/pfcpsim /bin
ENTRYPOINT [ "/bin/pfcpsim" ]
