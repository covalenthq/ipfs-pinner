# IPFS-Pinner

- [Introduction](#introduction)
- [Usage as a library](#usage-as-a-library)
  - [A note on CID and CAR files](#a-note-on-cid-and-car-files)
- [Running ipfs-pinner server](#running-ipfs-pinner-server)
  - [Upload a file](#upload-a-file)
  - [Download content (given cid)](#download-content-given-cid)
  - [Find the cid given some content](#find-the-cid-given-some-content)
- [Running ipfs-pinner server with docker](#running-ipfs-pinner-server-with-docker)
  - [Docker Volume setup](#docker-volume-setup)
  - [Port mapping setup](#port-mapping-setup)
- [Development](#development)
  - [Generate go http client go bindings via openapi](#generate-go-http-client-go-bindings-via-openapi)
- [Known Issues](#known-issues)
  - [Permission issue](#permission-issue)
  - [Netscan alert issue](#netscan-alert-issue)
  - [Updating IPFS http gateways](#updating-ipfs-http-gateways)

## Introduction

- A wrapper on top of ipfs node, utilising go-ipfs as a library.
- Extended support for custom file upload endpoints provided by pinata & web3.storage.
- Content archive file generation and lightweight deterministic CID generation on client side (using CARs).
- It can be used as a go library (see `binary/main.go` for usage) or as a http server.

## Usage as a library

Create a `PinnerNode` and start utilizing the services it is composed of. Check `binary/main.go` for detailed usage.

### A note on CID and CAR files

When you do a `ipfs add`, the content of the files are chunked and a merkle DAG is created which is used to compute the root CID. Now there are various parameters which can determine the structure of the DAG (see `ipfs help add`), and using different parameters results in different CIDs. As an example, web3.storage and pinata reported different CIDs for files greater than 256KB, because web3.storage chunked the data at 1 MB while pinata used the default value (256Kb).

To avoid this issue, the merkle DAG thus generated is exported into special files called content archives (CAR), and uploaded as CAR files. Thus, the merkle DAG structure is encapsulated in the files uploaded, and the same CID can now be used no matter what the pinning service is.

## Running ipfs-pinner server

1. Set the environment variable `WEB3_JWT`

2. to start a server which listens for request on 3001 port, run:

```bash
make clean server-dbg run
```

NOTE: If you want more control over CLI params, you can run the server binary (after `make clean server-dbg`):

```bash
./build/bin/server -jwt <WEB3_JWT> -port 3001
```

NOTE: If you get some error when running this, check if the diagnostic is there in [known issues](#known-issues)

ipfs-pinner can be run as a server and allows two functionalities currently - `/get` and `/upload`

### Upload a file

- Submit a request to upload a file (using multipart/form-data):

```bash
➜ curl -F "filedata=@file_to_upload" http://127.0.0.1:3001/upload
{"cid": "QmUqcL1RwbnbQ3FzmnT1aeRk8g8L5naKinJd5hCuPXxbZ2"}
```

Failures will be reported (via "error" field in the json). E.g:

```bash
➜ curl -F "filedata=@non_existent_file" http://127.0.0.1:3001/upload
{"error": "open not_exist_file: no such file or directory"}
```

There's a timeout (check code for value), on timeout the error message returned is:

```bash
{"error": "context deadline exceeded"}
```

### Download content (given cid)

- Download a file:

If the request succeeds the raw content is sent back and it can be outputted in a file using curl. e.g.

```bash
➜ curl -XGET http://127.0.0.1:3001/get\?cid\=bafybeifzst7cbujrqemiulznrkttouzshnqkrajiib5fp5te53ojs5sl5u --output file.jpeg
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  806k    0  806k    0     0  91.2M      0 --:--:-- --:--:-- --:--:--  262M
```

Now, if the data is present in local ipfs store, it'll be returned from there. Otherwise, it has to interact with other IPFS peers to find and fetch the data.
There's a timeout (check code for value) for the download request, if it doesn't succeed in that time, an error message is returned:

```bash
{"error": "context deadline exceeded"}
```

### Find the cid given some content

```bash
➜ curl -F "filedata=@LICENSE" http://127.0.0.1:3001/cid
{"cid": "bafkreicszve3ewhhrgobm366mdctki2m2qwzide5e54zh5aifnesg3ofne"}%
```

## Running ipfs-pinner server with docker

We can also run the ipfs-pinner server via docker.
for ipfs-pinner to function properly with docker, we need

- Docker volumes, to host the added data and persist it across container restarts/kill etc.
- Expose the ports that ipfs needs from the docker.

Docker run command should have:

- Volumes for data persistence
- Port mappings
- JWT token passed in the env

```bash
docker buildx create --name builder --use --platform=linux/amd64,linux/arm64  && docker buildx build --platform=linux/amd64,linux/arm64 . -t gcr.io/covalent-project/ipfs-pinner:latest
```

Now, we can run the container:

```bash
docker container run --detach --name ipfs-pinner-instance \
       --volume /tmp/data/.ipfs/:/root/.ipfs/  \
       -p 4001:4001 -p 3001:3001  \
       --env WEB3_JWT=$WEB3_JWT \
    <image-id>
```

### Docker Volume setup

There's 1 docker volume that needs to be shared (and persisted) between the container and the host - this `~/.ipfs` directory needs to have its lifecycle unaffected by container lifecycle (since it contains the merklelized nodes, blockstore etc.), and so that is docker volume managed.  

### Port mapping setup

:4001 : swarm port for p2p  
:8080 - http gateway (used by encapsulated ipfs-node)
:5001: local api (should be bound to 127.0.0.1 only, and must never be exposed publicly as it allows one to control the ipfs node; also used by encapsulated ipfs-node)  
:3001: The ipfs-pinner itself exposes its REST API on this port

<B> Out of the above, only the swarm port and the REST api port (3001) are essential.</B>  

---

## Development

### Generate go http client go bindings via openapi

- use `./generate_gobindings.sh` to generate the golang bindings (for pinning services of pinata and web3.storage).

- There are some fixes you would need to do (missing braces etc).

## Known Issues

### permission issue

If while using the ipfs-pinner as a server, you come across any permissions issues with logs such as

```log
Permission denied: Unable to access ./ipfs/plugins ...
etc
```

Or above fails with a message about permission issues to access  ~/.ipfs/*, run the following against the ipfs directory and try again.

```bash
sudo chmod -R 700 ~/.ipfs
```

### Netscan alert issue

If while using ipfs-pinner, a netscan alert is triggered due the exposed usage of port 4001 (swarm port for p2p) while ipfs tries to look for ipfs nodes in an internal network, this can be avoided by running ipfs as a server by updating the config in the following steps.

  1. Shut down the nodes using ipfs.
  2. Apply the config.
  3. Restart the nodes.

```bash
sudo systemctl stop bsp-agent.service
sudo systemctl stop ipfs-pinner.service
ipfs config profile apply server

{
"API": {"HTTPHeaders":{}},
"Addresses": {
    "API": "/ip4/127.0.0.1/tcp/5001",
    "Announce": [],
    "AppendAnnounce": [],
    "Gateway": "/ip4/127.0.0.1/tcp/8080",
    "NoAnnounce": {
        << "": "/ip4/10.0.0.0/ipcidr/8",
..
...
....
.....
    <> "DisableNatPortMap": false,
    ** "DisableNatPortMap": true,
    "RelayClient": {},
    "RelayService": {},
    "ResourceMgr": {},
    "Transports": {"Multiplexers":{},"Network":{},"Security":{}}
    }
sudo systemctl start ipfs-pinner.service
sudo systemctl start bsp-agent.service
```

This effectively disables local host discovery and is recommended when running IPFS on machines with public IPv4 addresses.

### Updating IPFS http gateways

ipfs-pinner currently uses some known IPFS gateways to fetch content. These gateways are expected to be run and maintained for a long time, but if you need to update the gateways list due to one of the going down, or a more efficient gateway being introduced etc. you can change the list:

```bash
./build/bin/server -jwt <WEB3_JWT> -port 3001 -ipfs-gateway-urls "https://w3s.link/ipfs/%s,https://dweb.link/ipfs/%s,https://ipfs.io/ipfs/%s"
```

The `-ipfs-gateways-urls` is a comma separated list of http urls with a `%s` present in it, which is formatted to replace the IPFS content identifier (CID) in it.
