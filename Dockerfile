# SPDX-License-Identifier: Apache-2.0
# Copyright 2022 Open Networking Foundation

# Stage pfcpsim-build: builds the pfcpsim docker image
FROM golang:alpine AS builder
WORKDIR /pfcpsimctl

COPY go.mod ./go.mod
COPY go.sum ./go.sum

RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -o /bin/pfcpsimctl cmd/pfcpsimctl/main.go
RUN CGO_ENABLED=0 go build -o /bin/pfcpsim cmd/pfcpsim/main.go

# Stage pfcpsimctl: runtime image of pfcpsimctl (client)
FROM golang:alpine AS pfcpsimctl

COPY --from=builder /bin/pfcpsimctl /bin
ENTRYPOINT [ "/bin/pfcpsimctl" ]

# Stage pfcpsim: runtime image of pfcpsim server
FROM golang:alpine AS pfcpsim

RUN apk update && apk add net-tools

COPY --from=builder /bin/pfcpsim /bin
ENTRYPOINT [ "/bin/pfcpsim" ]
