# Build - first phase
FROM golang:1.17-alpine as builder
RUN mkdir /build
WORKDIR /build
COPY . .
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o ipfs-server ./server/main.go

# Runtime -  second phase.
FROM alpine:3.15.0
RUN mkdir /app
WORKDIR /app
RUN apk update && apk add --no-cache bash=5.1.16-r0
COPY --from=builder /build/ipfs-server /app
SHELL ["/bin/bash", "-c"]
RUN chmod +x ./ipfs-server
ENTRYPOINT [ "/bin/bash", "-l", "-c" ]
CMD [ "./ipfs-server -port 3000 -jwt $WEB3_JWT" ]

EXPOSE 3000
# Swarm TCP; should be exposed to the public
EXPOSE 4001
# Daemon API; must not be exposed publicly but to client services under you control
EXPOSE 5001
# Web Gateway; can be exposed publicly with a proxy, e.g. as https://ipfs.example.org
EXPOSE 8080
# Swarm Websockets; must be exposed publicly when the node is listening using the websocket transport (/ipX/.../tcp/8081/ws).
EXPOSE 8081