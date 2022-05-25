# ipfs-pinner
- go client for IPFS pinning service api
- extended support for custom file upload endpoints provided by pinata & web3.storage
- car file generation and lightweight deterministic CID generation on client side (using cars).

### generate go http client go bindings via openapi
- use `./generate_gobindings.sh` to generate the golang bindings (for pinning services of pinata and web3.storage). 
- There are some fixes you would need to do (missing braces etc.)
