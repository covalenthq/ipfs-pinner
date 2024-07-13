# Build - first phase
FROM golang:1.22-alpine as builder
RUN mkdir /build
WORKDIR /build
COPY . .
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o ipfs-server ./server/main.go

# Runtime -  second phase.
FROM alpine:3.20
RUN mkdir /app
WORKDIR /app
RUN apk update && apk add --no-cache bash nodejs npm git && npm install -g @web3-storage/w3cli
COPY --from=builder --chmod=700 /build/ipfs-server /app

RUN apk del git && rm -rf /var/cache/apk/* /root/.npm /tmp/*

HEALTHCHECK --interval=10s --timeout=5s CMD wget --no-verbose --tries=1 --spider localhost:3001/health

ENTRYPOINT [ "/bin/bash", "-l", "-c" ]
CMD [ "./ipfs-server -port 3001 -w3-agent-key $W3_AGENT_KEY -w3-delegation-file $W3_DELEGATION_FILE --enable-gc" ]

# ipfs-pinner API server;
EXPOSE 3001
# Swarm TCP; should be exposed to the public
EXPOSE 4001
# Daemon API; must not be exposed publicly but to client services under you control
EXPOSE 5001
# Web Gateway; can be exposed publicly with a proxy, e.g. as https://ipfs.example.org
EXPOSE 8080
# Swarm Websockets; must be exposed publicly when the node is listening using the websocket transport (/ipX/.../tcp/8081/ws).
EXPOSE 8081