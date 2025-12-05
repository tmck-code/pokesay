FROM golang:1.24-alpine

WORKDIR /usr/local/src

RUN apk update \
    && apk upgrade \
    && apk add --no-cache jq tree bash dpkg

RUN adduser -D u

ADD go.* /usr/local/src/
ADD src/ /usr/local/src/src/

RUN go mod tidy \
    && go get -v github.com/mitchellh/go-wordwrap \
    && go install gotest.tools/gotestsum@latest \
    && chown -R u:u /go

ADD . .
