# SPDX-License-Identifier: Apache-2.0
# Copyright 2022 Open Networking Foundation

# Stage pfcpsimctl-build: builds the pfcpsimctl docker image
FROM golang AS pfcpsim-build
WORKDIR /pfcpsimctl

COPY go.mod ./go.mod
COPY go.sum ./go.sum

RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -o /bin/pfcpsimctl cmd/pfcpsimctl/main.go
RUN CGO_ENABLED=0 go build -o /bin/pfcpsim-server cmd/pfcpsim-server/main.go

# Stage pfcpsimctl: runtime image of pfcpsimctl (client)
FROM golang AS pfcpsim-client

COPY --from=pfcpsim-build /bin/pfcpsimctl /bin
ENTRYPOINT [ "/bin/pfcpsimctl" ]

# Stage pfcpsim: runtime image of pfcpsim server
FROM golang AS pfcpsim-server

RUN apt-get update && apt-get install -y net-tools

COPY --from=pfcpsim-build /bin/pfcpsim-server /bin
ENTRYPOINT [ "/bin/pfcpsim-server" ]
