[![Go Report Card](https://goreportcard.com/badge/github.com/omec-project/pfcpsim)](https://goreportcard.com/report/github.com/omec-project/pfcpsim)

# pfcpsim
pfcpsim is a simulator to interact with PFCP agents. Can be used to simulate a 4G SGW-C / 5G SMF.

## Overview
pfcpsim is designed to work within a containerized environment. The docker image comes with both client (`pfcpctl`) and server (`pfcpsim`).

`PFCPClient` is embedded in a gRPC Server. Interaction between pfcpsim and pfcpctl is performed through RPCs, as shown in the following schema: 

![Alt text](docs/images/schema.svg?raw=true "schema")

## Getting Started
#### 1. Build the container locally (You can skip this once docker image is published by the CI):
```bash
make build-pfcpsim
```

#### 2. Create the container. Use `-p` to set a custom gRPC listening port (default is 54321)
```bash
docker container run --rm -d --name pfcpsim pfcpsim:0.1.0-dev -p 12345
```

#### 3. Use `pfcpctl` to configure server's remote peer address and N3 interface address:
```bash
docker exec pfcpsim pfcpctl --server localhost:12345 -c configure --n3-addr <N3-interface-address> --remote-peer <PFCP-server-address>
```
 - `--server`: gRPC server address.
 - `-c`: command to execute.
 - `--n3-addr`: address of the N3 Interface between UPF and nodeB.
 - `--remote-peer`: address of the PFCP server. It supports the override of the IANA PFCP port (e.g. `10.0.0.1:8888`).

To list all the available commands just append `--help`, when executing `pfcpctl`.

#### 4. `associate` command will connect to remote peer set in the previous configuration step and perform an association.
```bash
docker exec pfcpsim pfcpctl --server localhost:12345 -c associate
```

#### 5. Create 5 sessions
```bash
docker exec pfcpsim pfcpctl --server localhost:12345 -c create --count 5 --baseID 2 --ue-pool <CIDR-IP-pool> --nb-addr <NodeB-address>
```
 - `--count` the amount of sessions to create
 - `--baseID` the base ID used to incrementally create sessions
 - `--ue-pool` the IP pool from which UE addresses will be generated (e.g. `17.0.0.0/24`)
 - `--nb-addr` the nodeB address 

#### 6. Delete the sessions
```bash
docker exec pfcpsim pfcpctl --server localhost:12345 -c delete --count 5 --baseID 2
```

#### 7. `disassociate` command will perform disassociation and close connection with remote peer.
```bash
docker exec pfcpsim pfcpctl --server localhost:12345 -c disassociate
```
