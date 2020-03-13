FROM golang:1.13 as build

ARG VERSION

ADD . /go/src/github.com/leeif/go-web-template

WORKDIR /go/src/github.com/leeif/go-web-template

RUN  export GO111MODULE=on GOPROXY=https://proxy.golang.org && \ 
  go build -ldflags="-X 'main.VERSION=${VERSION}'" -o server cmd/server/main.go && \
  go build -o db-migrate cmd/db-migrate/main.go

FROM ubuntu:18.04

RUN apt-get update && apt-get install -y wget

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz

ENV ConfigFile /etc/go-web-template/config.json

COPY --from=build /go/src/github.com/leeif/go-web-template/server /usr/bin/

COPY --from=build /go/src/github.com/leeif/go-web-template/db-migrate /usr/bin/

ENTRYPOINT ["/usr/bin/server"]