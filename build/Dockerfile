FROM golang:1.24

WORKDIR /usr/local/src

ENV DEBIAN_FRONTEND=noninteractive
ENV BUILD_DEPS="make gcc"
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        $BUILD_DEPS libmagickwand-dev ncurses-dev jq imagemagick tree \
    && git clone -q --depth 1 https://github.com/denilsonsa/img2xterm \
    && (cd img2xterm && make && make install) \
    && rm -rf img2xterm \
    && apt-get purge -y $BUILD_DEPS \
    && apt-get autoremove --purge -y \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /tmp/* /var/tmp/*


RUN useradd u -m

WORKDIR /usr/local/src/
ADD go.* /usr/local/src/
ADD src/ /usr/local/src/src/
RUN go mod tidy \
    && go get -v github.com/mitchellh/go-wordwrap \
    && go install gotest.tools/gotestsum@latest \
    && chown -R u:u /go

ADD . .
