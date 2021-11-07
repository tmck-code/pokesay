ARG GOOS=linux
ARG GOARCH=amd64

FROM golang:latest

WORKDIR /usr/local/src

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        parallel make gcc libmagickwand-dev ncurses-dev fortune-mod fortunes

RUN git clone --depth 1 git://github.com/msikma/pokesprite && \
    cp -R pokesprite/icons/ . && \
    rm -rf pokesprite

ENV PATH "/usr/lib/x86_64-linux-gnu/ImageMagick-6.9.11/bin-q16:/usr/include/ImageMagick-6/:$PATH"
RUN git clone --depth 1 git://github.com/rossy/img2xterm && \
    (cd img2xterm && make && make install)

ENV GOOS $GOOS
ENV GOARCH $GOARCH

WORKDIR /usr/local/src/
RUN go mod init github.com/tmck-code/pokesay-go \
    && go install github.com/go-bindata/go-bindata/...@latest \
    && go get github.com/mitchellh/go-wordwrap

ADD build/build_cows.sh ./

ENV DOCKER_BUILD_DIR  /usr/local/src
ENV DOCKER_OUTPUT_DIR /tmp

ADD *.go ./

RUN ./build_cows.sh \
    && go-bindata /tmp/cows/... \
    && go build pokesay.go bindata.go

RUN /usr/games/fortune | ./pokesay
