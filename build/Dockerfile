FROM golang:latest

WORKDIR /usr/local/src

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        parallel make gcc libmagickwand-dev ncurses-dev fortune-mod fortunes

RUN git clone --depth 1 git://github.com/denilsonsa/img2xterm \
    && (cd img2xterm && make && make install)

ENV DOCKER_BUILD_DIR /tmp/original
ENV DOCKER_OUTPUT_DIR /tmp/cows

RUN mkdir -p /tmp/original /tmp/cows
WORKDIR /tmp/original

RUN git clone --depth 1 git://github.com/msikma/pokesprite

WORKDIR /usr/local/src/
ADD build/build_cows.go build/bindata.sh /usr/local/src/build/
ADD go.* /usr/local/src/
RUN go mod tidy \
    && go get -v github.com/mitchellh/go-wordwrap \
    && go install github.com/go-bindata/go-bindata/...@latest

RUN go run /usr/local/src/build/build_cows.go \
        -from /tmp/original/pokesprite/ \
        -to /tmp/cows/ \
        -skip '["resources/", "misc/"]'

ADD cmd/ /usr/local/src/cmd/
RUN go-bindata -nometadata -nomemcopy -prefix /tmp/ /tmp/cows/... && \
        ./build/bindata.sh


RUN echo "building darwin / amd64"  && GOOS=darwin  GOARCH=amd64 go build -o pokesay-darwin-amd64  cmd/pokesay.go cmd/bindata.go && \
    echo "building linux / amd64"   && GOOS=linux   GOARCH=amd64 go build -o pokesay-linux-amd64   cmd/pokesay.go cmd/bindata.go && \
    echo "building windows / amd64" && GOOS=windows GOARCH=amd64 go build -o pokesay-windows-amd64.exe cmd/pokesay.go cmd/bindata.go && \
    echo "building android / arm64" && GOOS=android GOARCH=arm64 go build -o pokesay-android-arm64 cmd/pokesay.go cmd/bindata.go
