FROM golang:1.17-alpine
COPY . /usr/src/app
COPY entry.sh /usr/src/app

WORKDIR /usr/src/app
RUN go mod download

RUN cd server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

RUN apk update && apk add --no-cache bash
SHELL ["/bin/bash", "-c"]
RUN chmod +x /usr/src/app/entry.sh
RUN chmod +x /usr/src/app/server
ENTRYPOINT [ "/bin/bash", "-l", "-c" ]
CMD [ "./entry.sh" ]
EXPOSE 3001:3001