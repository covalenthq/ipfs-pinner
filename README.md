# ipfs-pinner
- go client for IPFS pinning service api
- extended support for custom file upload endpoints provided by pinata & web3.storage
- car file generation and lightweight deterministic CID generation on client side (using cars).

### generate go http client go bindings via openapi
- use `./generate_gobindings.sh` to generate the golang bindings (for pinning services of pinata and web3.storage). 
- There are some fixes you would need to do (missing braces etc.)

### Run the server
to start a server which listens for request on a particular port, run:
```bash
go run main.go -port 3000 -jwt "<jwt_token>"
```

- submit a request to upload a file:
```bash
➜ curl -XGET http://127.0.0.1:3000/pin\?filePath\=/Users/sudeep/repos/elixir-koans/mix.exs
{"cid": "QmUqcL1RwbnbQ3FzmnT1aeRk8g8L5naKinJd5hCuPXxbZ2"}
```

- failures will be reported (via "error" field in the json). E.g:
```bash
➜ curl -XGET http://127.0.0.1:3000/pin\?filePath\=not_exist_file
{"error": "open not_exist_file: no such file or directory"}
```
