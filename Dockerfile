# SPDX-License-Identifier: Apache-2.0
# Copyright 2022 Open Networking Foundation

# Stage pfcpsimctl-build: builds the pfcpsimctl docker image
# FIXME check all paths
FROM golang AS pfcpsim-client-build
WORKDIR /pfcpsim-client

COPY go.mod ./go.mod
COPY go.sum ./go.sum

RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -o /bin/pfcpsimctl cmd/pfcpsimctl/main.go

# Stage pfcpsimctl: runtime image of pfcpsimctl (client)
FROM golang AS pfcpsim-client

COPY --from=pfcpsim-client-build /bin/pfcpsim-client /bin
ENTRYPOINT [ "/bin/pfcpsimctl" ]

# Stage pfcpsim: runtime image of pfcpsim server
FROM golang AS pfcpsim-client

RUN apt-get update && apt-get install -y net-tools iputils-ping

COPY --from=pfcpsim-client-build /bin/pfcpsim-client /bin
ENTRYPOINT [ "/bin/pfcpsim-server" ]
