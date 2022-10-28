# ipfs-pinner
- go client for IPFS pinning service api
- extended support for custom file upload endpoints provided by pinata & web3.storage
- car file generation and lightweight deterministic CID generation on client side (using cars).
- it can be used as a go library (see `binary/main.go` for usage) or as a http server

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

there's a timeout (check code for value), on timeout the error message returned is:
```bash
{"error": "context deadline exceeded"}
```

- download a file
if the request succeeds the raw content is sent back and it can be outputted in a file using curl. e.g.
```bash
➜ curl -XGET http://127.0.0.1:3000/get\?cid\=bafybeifzst7cbujrqemiulznrkttouzshnqkrajiib5fp5te53ojs5sl5u --output file.jpeg
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  806k    0  806k    0     0  91.2M      0 --:--:-- --:--:-- --:--:--  262M
```

Now, if the data is present in local ipfs store, it'll be returned from there. Otherwise, it has to interact with other IPFS peers to find 
and fetch the data.
There's a timeout (check code for value) for the download request, if it doesn't succeed in that time, the error message returned is:
```bash
{"error": "context deadline exceeded"}
```
