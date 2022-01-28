# SPDX-License-Identifier: Apache-2.0
# Copyright 2022 Open Networking Foundation

# Stage pfcpsim-client-build: builds the pfcpsim-client docker image
FROM golang AS pfcpsim-client-build
WORKDIR /pfcpsim-client

COPY go.mod ./go.mod
COPY go.sum ./go.sum

RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -o /bin/pfcpsim-client cmd/pfcpsim-client/main.go

# Stage pfcpsim-client: runtime image of pfcpsim-client
FROM golang AS pfcpsim-client

RUN apt-get update && apt-get install -y net-tools iputils-ping

COPY --from=pfcpsim-client-build /bin/pfcpsim-client /bin
ENTRYPOINT [ "/bin/pfcpsim-client" ]
