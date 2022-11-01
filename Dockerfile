FROM golang:1.17-alpine

COPY . /usr/src/app/repo

WORKDIR /usr/src/app
RUN cd repo && go mod download

RUN cd repo/server && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../server .
RUN cd repo/server && go clean -modcache
RUN rm -r /usr/src/app/repo

RUN apk update && apk add --no-cache bash

SHELL ["/bin/bash", "-c"]

RUN chmod +x ./server

ENTRYPOINT [ "/bin/bash", "-l", "-c" ]
CMD [ "./server -port 3000 -jwt $WEB3_JWT" ]

EXPOSE 3000
# Swarm TCP; should be exposed to the public
EXPOSE 4001
# Daemon API; must not be exposed publicly but to client services under you control
EXPOSE 5001
# Web Gateway; can be exposed publicly with a proxy, e.g. as https://ipfs.example.org
EXPOSE 8080
# Swarm Websockets; must be exposed publicly when the node is listening using the websocket transport (/ipX/.../tcp/8081/ws).
EXPOSE 8081