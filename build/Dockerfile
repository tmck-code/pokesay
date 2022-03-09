FROM golang:latest

WORKDIR /usr/local/src

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        parallel make gcc libmagickwand-dev ncurses-dev fortune-mod fortunes

ENV PATH "/usr/lib/x86_64-linux-gnu/ImageMagick-6.9.11/bin-q16:/usr/include/ImageMagick-6/:$PATH"
RUN git clone --depth 1 git://github.com/rossy/img2xterm \
    && (cd img2xterm && make && make install)

ENV DOCKER_BUILD_DIR /tmp/original
ENV DOCKER_OUTPUT_DIR /tmp/converted
ADD build/build_cows.sh /usr/local/src/build/

RUN mkdir -p /tmp/original /tmp/converted
WORKDIR /tmp/original

RUN git clone --depth 1 git://github.com/msikma/pokesprite \
    && cp -R pokesprite/icons/ . \
    && rm -rf pokesprite \
    && /usr/local/src/build/build_cows.sh

WORKDIR /usr/local/src/
ADD go.* /usr/local/src/
RUN go mod tidy \
    && go get -v github.com/mitchellh/go-wordwrap \
    && go get google.golang.org/protobuf/proto

ADD cmd/ /usr/local/src/cmd/

# RUN protoc --proto_path=cmd --go_out=cmd/ --go_opt=paths=source_relative cmd/pokemon.proto
RUN echo "building darwin / amd64"  && GOOS=darwin  GOARCH=amd64 go build -o pokesay-darwin-amd64  cmd/pokesay.go cmd/pokemon.pb.go && \
    echo "building linux / amd64"   && GOOS=linux   GOARCH=amd64 go build -o pokesay-linux-amd64   cmd/pokesay.go cmd/pokemon.pb.go && \
    echo "building windows / amd64" && GOOS=windows GOARCH=amd64 go build -o pokesay-windows-amd64 cmd/pokesay.go cmd/pokemon.pb.go && \
    echo "building android / arm64" && GOOS=android GOARCH=arm64 go build -o pokesay-android-arm64 cmd/pokesay.go cmd/pokemon.pb.go
