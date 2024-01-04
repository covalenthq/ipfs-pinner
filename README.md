# IPFS-Pinner

- [IPFS-Pinner](#ipfs-pinner)
  - [Introduction](#introduction)
  - [Usage as a library](#usage-as-a-library)
    - [A note on CID and CAR files](#a-note-on-cid-and-car-files)
  - [Running ipfs-pinner server](#running-ipfs-pinner-server)
    - [Upload a file](#upload-a-file)
    - [Download content (given cid)](#download-content-given-cid)
    - [Find the cid given some content](#find-the-cid-given-some-content)
  - [migration to UCAN and capabilities setup](#migration-to-ucan-and-capabilities-setup)
    - [setting up w3cli](#setting-up-w3cli)
      - [installation](#installation)
      - [login and check spaces](#login-and-check-spaces)
      - [generate ucan key](#generate-ucan-key)
      - [create delegation to store/add and upload/add](#create-delegation-to-storeadd-and-uploadadd)
      - [communicate to the operator](#communicate-to-the-operator)
      - [operator invocation](#operator-invocation)
  - [Running ipfs-pinner server with docker](#running-ipfs-pinner-server-with-docker)
    - [Docker Volume setup](#docker-volume-setup)
    - [Port mapping setup](#port-mapping-setup)
  - [Development](#development)
    - [Generate go http client go bindings via openapi](#generate-go-http-client-go-bindings-via-openapi)
    - [Improvements](#improvements)
  - [Known Issues](#known-issues)
    - [Permission Issue](#permission-issue)
    - [UDP buffer size warning](#udp-buffer-size-warning)
    - [Netscan alert issue](#netscan-alert-issue)
    - [Using a different directory than ~/.ipfs](#using-a-different-directory-than-ipfs)
    - [Updating IPFS http gateways](#updating-ipfs-http-gateways)

## Introduction

- A wrapper on top of ipfs node, utilising go-ipfs as a library.
- Extended support for custom file upload endpoints provided by web3.storage.
- Content archive file generation and lightweight deterministic CID generation on client side (using CARs).
- It can be used as a go library (see `binary/main.go` for usage) or as a http server.

## Usage as a library

Create a `PinnerNode` and start utilizing the services it is composed of. Check `binary/main.go` for detailed usage.

### A note on CID and CAR files

When you do a `ipfs add`, the content of the files are chunked and a merkle DAG is created which is used to compute the root CID. Now there are various parameters which can determine the structure of the DAG (see `ipfs help add`), and using different parameters results in different CIDs. As an example, web3.storage and pinata reported different CIDs for files greater than 256KB, because web3.storage chunked the data at 1 MB while pinata used the default value (256Kb).

To avoid this issue, the merkle DAG thus generated is exported into special files called content archives (CAR), and uploaded as CAR files. Thus, the merkle DAG structure is encapsulated in the files uploaded, and the same CID can now be used no matter what the pinning service is.

## Running ipfs-pinner server

1. Get the agent key, did and delegation proof from Covalent

2. build the server and run:

```bash
make clean server-dbg
```

NOTE: If you want more control over CLI params, you can run the server binary (after `make clean server-dbg`):

```bash
./build/bin/server -jwt <WEB3_JWT> -port 3001
./build/bin/server -w3-agent-key <AGENT_KEY> -w3-delegation-file <DELEGATION_PROOF_FILE_PATH>
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


## migration to UCAN and capabilities setup
- web3.storage is sunsetting its custom upload endpoint (on 9th January, 2024), and we need to migrate from using that to w3up.
- w3up uses UCAN which is a capabilities-based authorization system (learn more [here](https://web3.storage/docs/concepts/ucans-and-web3storage/)). 
- In this setup, the "central account" (owned by Covalent) sets up a "space" (think namespace for data). The central account (controlled by the email) is delegated the capabilty to operate on this space.
- among other capabilties, the central account can delegate certain capabilities (like uploading to space) to other **agents**. This has to be done at our end, and scripts will be made available for it in this repo.
- once an agent is granted the capability, we share the credentials with the operators, who run ipfs-pinner with it, and can then upload or fetch.


### setting up w3cli

- Create a web3.storage account in the [console](https://console.web3.storage/). 
- Create a space which you want to use to upload artifacts. We want to use different spaces for different artifacts to keep a clear separation.

We'll use w3cli to login and create a new space and register.

#### installation
```bash
➜ npm install -g @web3-storage/w3cli

➜ w3 --version
w3, 7.0.3
```

#### login and check spaces
```bash
➜ w3 login sudeep@covalenthq.com

➜ w3 space ls
  did:key:z6MkgSK6VEu3bvrAFtYNyjsnzG7dVXzYi3yT5TasEgeaQrCe mock_artifacts

➜ w3 space use did:key:z6MkgSK6VEu3bvrAFtYNyjsnzG7dVXzYi3yT5TasEgeaQrCe
did:key:z6MkgSK6VEu3bvrAFtYNyjsnzG7dVXzYi3yT5TasEgeaQrCe
```

The did key is the identifier for this space. Now let's generate some DIDs for an operator and delegate upload capabilities to it.

#### generate ucan key
```bash
➜ npx ucan-key ed --json
{
  "did": "did:key:z6MkpzWw1fDZYMpESgVKFAT87SZAuHiCQZVBC3hmQjB18Nzj",
  "key": "MgCbc48J8n+BMdzA4XxwYOaKmdu5Ov34jE71U8vV07IVIjO0BnJa05mNMcB8GSz1lib014QAhvAxorG6zACrstm6PBGA="
}
```

#### create delegation to store/add and upload/add

```bash
➜ w3 delegation create -c 'store/add' -c 'upload/add' did:key:z6MkpzWw1fDZYMpESgVKFAT87SZAuHiCQZVBC3hmQjB18Nzj -o proof.out
```


Copy the output. This is the delegation string.

#### communicate to the operator

Provide the operator with the `did`, `key` string + `proof.out` file. These will be passed to operator's setup of the 
ipfs-pinner, which can then make the delegations.


#### operator invocation

the operator can pass the `key` for `-w3-agent-key` and proof file in `-w3-delegation-file` flag.

```bash
go run server/main.go -w3-agent-key <agent-key> -w3-delegation-file ./proof.out
ipfs-pinner
ipfs-pinner Version: 0.1.16
Architecture: arm64
Go Version: go1.20.5
Operating System: darwin
GOPATH=/Users/sudeep/go/
GOROOT=/usr/local/go
2024/01/04 15:52:05 agent did: did:key:z6MkoLvhaiE9NRYs3vJcynCM8CeyP8hXduWhE5Ter2U2x93y
generating 2048-bit RSA keypair...done
peer identity: QmY49BMJdGneQjJAbTPrGSqaQcLjpCE1WFkRBP6XZEHd6i
2024/01/04 15:52:09 setting up w3up for uploads....
2024/01/04 15:52:10 w3up agent did: did:key:z6MkoLvhaiE9NRYs3vJcynCM8CeyP8hXduWhE5Ter2U2x93y
2024/01/04 15:52:10 w3up space did: did:key:z6MkgSK6VEu3bvrAFtYNyjsnzG7dVXzYi3yT5TasEgeaQrCe
2024/01/04 15:52:10 w3up setup complete
2024/01/04 15:52:10 Listening...
2024/01/04 15:52:15 generated dag has root cid: bafybeigvijf76lcsjwcmkr6rmzovoiiqdog3muqs5vnplvf4jxh47shfiu
2024/01/04 15:52:15 car file location: /var/folders/w0/bf3y1c7d6ys15tq97ffk5qhw0000gn/T/3475885728.car
2024/01/04 15:53:06 w3 up output: {"root":{"/":"bafybeigvijf76lcsjwcmkr6rmzovoiiqdog3muqs5vnplvf4jxh47shfiu"}}
2024/01/04 15:53:28 uploaded file has root cid: bafybeigvijf76lcsjwcmkr6rmzovoiiqdog3muqs5vnplvf4jxh47shfiu
```



## Running ipfs-pinner server with docker

We can also run the ipfs-pinner server via docker.
for ipfs-pinner to function properly with docker, we need

- Docker volumes, to host the added data and persist it across container restarts/kill etc.
- Expose the ports that ipfs needs from the docker.

Docker run command should have:

- Volumes for data persistence; volumes containing the delegation proof
- Port mappings
- W3up agent key passed in the env

```bash
docker buildx create --name builder --use --platform=linux/amd64,linux/arm64  && docker buildx build --platform=linux/amd64,linux/arm64 . -t gcr.io/covalent-project/ipfs-pinner:latest
```

Now, we can run the container:

```bash
docker container run --detach --name ipfs-pinner-instance \
       --volume /tmp/data/.ipfs/:/root/.ipfs/  \
       --volume /Users/sudeep/repos/ipfs/w3up-testing/w3up_proof:/root/w3up_proof/ \
       -p 3001:3001  \
       --env W3_AGENT_KEY=$W3_AGENT_KEY \
       --env W3_DELEGATION_FILE=/root/w3up_proof/proof.out
    <image-id>
```

### Docker Volume setup

There's 1 docker volume that needs to be shared (and persisted) between the container and the host - this `~/.ipfs` directory needs to have its lifecycle unaffected by container lifecycle (since it contains the merklelized nodes, blockstore etc.), and so that is docker volume managed.  

### Port mapping setup

:4001 : swarm port for p2p  (currently disabled)
:8080 - http gateway (used by encapsulated ipfs-node)
:5001: local api (should be bound to 127.0.0.1 only, and must never be exposed publicly as it allows one to control the ipfs node; also used by encapsulated ipfs-node)  
:3001: The ipfs-pinner itself exposes its REST API on this port

<B> Out of the above, only the swarm port and the REST api port (3001) are essential.</B>  

---

## Development

### Generate go http client go bindings via openapi

- use `./generate_gobindings.sh` to generate the golang bindings (for pinning services of pinata and web3.storage).

- There are some fixes you would need to do (missing braces etc).

### Improvements

- remove the pinning service yaml (none of the pinning service api is currently used). Directly use the web3.storage goclient - https://github.com/web3-storage/go-w3s-client


## Known Issues

### Permission Issue

If while using the ipfs-pinner as a server, you come across any permissions issues with logs such as

```log
Permission denied: Unable to access ./ipfs/plugins ...
etc
```

Or above fails with a message about permission issues to access  ~/.ipfs/*, run the following against the ipfs directory and try again.

```bash
sudo chmod -R 700 ~/.ipfs
```

### UDP buffer size warning

On the start of ipfs-pinner, you might notice logs regarding UDP buffer size:
```
2023/07/26 05:43:14 failed to sufficiently increase receive buffer size (was: 208 kiB, wanted: 2048 kiB, got: 416 kiB). See https://github.com/lucas-clemente/quic-go/wiki/UDP-Receive-Buffer-Size for details.
```

Do go to the link mention in the log, or to https://github.com/quic-go/quic-go/wiki/UDP-Buffer-Sizes, it'll help QUIC function more effectively over high-bandwidth functions, reducing timeouts etc.

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

### Using a different directory than ~/.ipfs
`~/.ipfs` is the default location for the storage of local 'add'-ed contents. When a lot of content is added to the node, this directory can bloat up and use a lot of storage.  
Users would sometimes want to maintain a different volume to fulfil large storage requirements. As a solution, one can move the `.ipfs/` folder to a larger partition and **symlink** it on $HOME instead. Stop the ipfs-pinner before performing this process, and then restart for the effects to take place.

### Updating IPFS http gateways

ipfs-pinner currently uses some known IPFS gateways to fetch content. These gateways are expected to be run and maintained for a long time, but if you need to update the gateways list due to one of the going down, or a more efficient gateway being introduced etc. you can change the list:

```bash
./build/bin/server -ipfs-gateway-urls "https://w3s.link/ipfs/%s,https://dweb.link/ipfs/%s,https://ipfs.io/ipfs/%s" ##OTHER PARAMS
```

The `-ipfs-gateways-urls` is a comma separated list of http urls with a `%s` present in it, which is formatted to replace the IPFS content identifier (CID) in it.
