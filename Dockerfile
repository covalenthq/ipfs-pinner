FROM golang:1.17-alpine
COPY . /usr/src/app/temp

COPY entry.sh /usr/src/app
COPY ./server/temp.txt /usr/src/app

WORKDIR /usr/src/app
RUN cd temp && go mod download

RUN cd temp/server && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../server .

RUN apk update && apk add --no-cache bash

SHELL ["/bin/bash", "-c"]
EXPOSE 3001

RUN chmod +x entry.sh
RUN chmod +x ./server

RUN rm -r /usr/src/app/temp
ENTRYPOINT [ "/bin/bash", "-l", "-c" ]
CMD [ "./entry.sh" ]
EXPOSE 3001:3001