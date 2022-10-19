FROM golang:1.17-alpine
COPY . /usr/src/app
COPY entry.sh /usr/src/app

WORKDIR /usr/src/app
RUN go mod download
RUN mkdir built
RUN cd server && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../built/server .

RUN apk update && apk add --no-cache bash

SHELL ["/bin/bash", "-c"]
EXPOSE 3001:3001

RUN chmod +x entry.sh
RUN chmod +x ./built/server

ENTRYPOINT [ "/bin/bash", "-l", "-c" ]
CMD [ "./entry.sh" ]
EXPOSE 3001:3001