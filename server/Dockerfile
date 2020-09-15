FROM golang:latest

ENV GO111MODULE=on

MAINTAINER JTF "jtf.pub@gmail.com"
WORKDIR $GOPATH/src/go-docker-tutorial/server
ADD . $GOPATH/src/go-docker-tutorial/server
RUN GOPROXY="https://goproxy.io,direct" go get -d -v ./...
RUN GOPROXY="https://goproxy.io,direct" go build .
EXPOSE 9090
ENTRYPOINT ["./server"]