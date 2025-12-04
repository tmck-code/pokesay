FROM golang:1.24

WORKDIR /usr/local/src

RUN useradd u -m

WORKDIR /usr/local/src/
ADD go.* /usr/local/src/

RUN go mod tidy
ADD src/ /usr/local/src/src/

RUN go get -v github.com/mitchellh/go-wordwrap \
    && go install gotest.tools/gotestsum@latest \
    && chown -R u:u /go

ADD . .
