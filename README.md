# pfcpsim
pfcpsim is a simulator to interact with PFCP agents. Can be used to emulate a 4G SGW / 5G UPF.

## Interaction
pfcpsim is designed to work within a containerized environment. The docker image comes with both client (`pcpfsimctl`) and server (`pfcpsim`).

`PFCPClient` is embedded in a gRPC Server. Interaction between pfcpsim and pfcpsimctl is performed through RPC calls, as shown in the following schema: 

![Alt text](docs/images/schema.svg?raw=true "schema")

## Example
In the following example we will borrow `mock-up4` environment, from [upf-epc repository](https://github.com/omec-project/upf-epc/tree/master/test/integration).
Using mock-up4's `docker-compose.yml` will create a pfcp agent container listening on `0.0.0.0:8805`:

1. Build the container locally (You can skip this once docker image is published by the CI):
```bash
make build-pfcpsim
```
2. Create the container using host network (to be able to connect to 127.0.0.1:8805)
```bash
docker container run --rm -d --name pfcpsim --net host pfcpsim:0.1.0-dev
```

3. Use pfcpsimctl by using `docker exec` on the same container to configure the server (not passing any flag will use mock-up4's default values)
```bash
docker exec pfcpsim /bin/pfcpsimctl -c configure
```

4. Associate
```bash
docker exec pfcpsim /bin/pfcpsimctl -c associate
```

5. Create 5 sessions
```bash
docker exec pfcpsim /bin/pfcpsimctl -c create --count 5 --baseID 2 --ue-pool 17.0.0.0/24 --gnb-addr 198.18.0.10
```

6. Delete the sessions
```bash
docker exec pfcpsim /bin/pfcpsimctl -c delete --count 5 --baseID 2
```

7. Disconnect from remote peer
```bash
docker exec pfcpsim /bin/pfcpsimctl -c interrupt
```
