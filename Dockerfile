# Build - first phase. 
FROM golang:1.17-alpine as builder
COPY . /usr/src/app
COPY entry.sh /usr/local/bin

WORKDIR /usr/src/app
RUN go mod download

RUN cd server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

RUN apk update && apk add --no-cache bash=5.1.16-r0

SHELL ["/bin/bash", "-c"]
EXPOSE 3001:3001

RUN chmod +x entry.sh

ENTRYPOINT [ "/bin/bash", "-l", "-c" ]
CMD [ "./entry.sh" ]