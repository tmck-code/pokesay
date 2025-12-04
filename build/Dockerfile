FROM golang:1.24

WORKDIR /usr/local/src

ENV DEBIAN_FRONTEND=noninteractive
RUN apt update \
    && apt upgrade -y \
    && apt install -y --no-install-recommends jq tree

RUN useradd u -m

ADD go.* /usr/local/src/
ADD src/ /usr/local/src/src/

RUN go mod tidy \
    && go get -v github.com/mitchellh/go-wordwrap \
    && go install gotest.tools/gotestsum@latest \
    && chown -R u:u /go

ADD . .
