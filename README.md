[![Go Report Card](https://goreportcard.com/badge/github.com/omec-project/pfcpsim)](https://goreportcard.com/report/github.com/omec-project/pfcpsim)

# pfcpsim
pfcpsim is a simulator to interact with PFCP agents. Can be used to simulate a 4G SGW-C / 5G SMF.

## Overview

pfcpsim is designed to work within a containerized environment. The docker image comes with both client (`pfcpctl`) and server (`pfcpsim`).

`PFCPClient` is embedded in a gRPC Server. Interaction between pfcpsim and pfcpctl is performed through RPCs, as shown in the following schema: 

![Alt text](docs/images/schema.svg)

## Getting Started

#### 1. Create the container. Images are available on [Dockerhub](https://hub.docker.com/r/opennetworking/pfcpsim/tags):
```bash
docker container run --rm -d --name pfcpsim pfcpsim:<image_tag> -p 12345 --interface <interface-name>
```
 - `-p` (**optional**, default is 54321): to set a custom gRPC listening port
 - `--interface` (**optional**, default is first non-loopback interface): to specify a specific interface from which retrieve local IP address

#### 2. Use `pfcpctl` to configure server's remote peer address and N3 interface address:
```bash
docker exec pfcpsim pfcpctl --server localhost:12345 -c configure --n3-addr <N3-interface-address> --remote-peer <PFCP-server-address>
```
 - `--server`: (**optional**, Default is localhost:54321) gRPC server address.
 - `-c`: command to execute.
 - `--n3-addr`: address of the N3 Interface between UPF and nodeB.
 - `--remote-peer`: address of the PFCP server. It supports the override of the IANA PFCP port (e.g. `10.0.0.1:8888`).

To list all the available commands just append `--help`, when executing `pfcpctl`.

#### 3. `associate` command will connect to remote peer set in the previous configuration step and perform an association.
```bash
docker exec pfcpsim pfcpctl --server localhost:12345 -c associate
```

#### 4. Create 5 sessions
```bash
docker exec pfcpsim pfcpctl --server localhost:12345 -c create --count 5 --baseID 2 --ue-pool <CIDR-IP-pool> --nb-addr <NodeB-address>
```
 - `--count` the amount of sessions to create
 - `--baseID` the base ID used to incrementally create sessions
 - `--ue-pool` the IP pool from which UE addresses will be generated (e.g. `17.0.0.0/24`)
 - `--nb-addr` the nodeB address 

#### 5. Delete the sessions
```bash
docker exec pfcpsim pfcpctl --server localhost:12345 -c delete --count 5 --baseID 2
```

#### 6. `disassociate` command will perform disassociation and close connection with remote peer.
```bash
docker exec pfcpsim pfcpctl --server localhost:12345 -c disassociate
```
