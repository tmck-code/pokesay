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

WORKDIR /usr/local/src/
ADD . .
RUN go mod tidy \
    && go install -v github.com/kevinburke/go-bindata/...@latest \
    && go get github.com/mitchellh/go-wordwrap

ENV DOCKER_BUILD_DIR  /usr/local/src
ENV DOCKER_OUTPUT_DIR /tmp
RUN ./build/build_cows.sh \
    && go-bindata -nometadata -o cmd/bindata.go /tmp/cows/...

ARG GOOS=linux
ARG GOARCH=amd64
ENV GOOS $GOOS
ENV GOARCH $GOARCH

RUN echo "building $GOOS / $GOARCH" \
    && go build cmd/pokesay.go cmd/bindata.go
