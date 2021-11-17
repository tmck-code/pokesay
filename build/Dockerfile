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
ADD build/ /usr/local/src/build/

RUN mkdir -p /tmp/original /tmp/converted
WORKDIR /tmp/original

RUN git clone --depth 1 git://github.com/msikma/pokesprite \
    && cp -R pokesprite/icons/ . \
    && rm -rf pokesprite \
    && /usr/local/src/build/build_cows.sh

WORKDIR /usr/local/src/
ADD go.* /usr/local/src/
RUN go mod tidy \
    && go get -v github.com/mitchellh/go-wordwrap


ARG GOOS=linux
ARG GOARCH=amd64
ENV GOOS $GOOS
ENV GOARCH $GOARCH

ADD cmd/ /usr/local/src/cmd/
RUN echo "building $GOOS / $GOARCH" \
    && go build cmd/pokesay.go cmd/bindata.go

ADD . .