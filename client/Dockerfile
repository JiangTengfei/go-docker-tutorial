FROM golang:latest

ENV GO111MODULE=on

MAINTAINER JTF "jtf.pub@gmail.com"

WORKDIR $GOPATH/src/go-docker-tutorial/client

COPY . $GOPATH/src/go-docker-tutorial/client

RUN GOPROXY="https://goproxy.io,direct" go get -d -v ./...
RUN GOPROXY="https://goproxy.io,direct" go build .
EXPOSE 8080
ENTRYPOINT ["./client"]