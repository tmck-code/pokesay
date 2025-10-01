FROM golang:1.24

WORKDIR /usr/local/src

ENV DEBIAN_FRONTEND=noninteractive
ENV BUILD_DEPS="make gcc"
RUN apt-get update \
    && apt-get install -y --no-install-recommends $BUILD_DEPS libmagickwand-dev ncurses-dev jq imagemagick

RUN git clone -q --depth 1 https://github.com/denilsonsa/img2xterm \
    && (cd img2xterm && make && make install) \
    && rm -rf img2xterm \
    && apt-get purge -y $BUILD_DEPS \
    && apt-get autoremove --purge -y \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /tmp/* /var/tmp/*

WORKDIR /usr/local/src/
ADD src/ /usr/local/src/src/
ADD go.* /usr/local/src/
RUN go mod tidy \
    && go get -v github.com/mitchellh/go-wordwrap \
    && go install gotest.tools/gotestsum@latest

RUN useradd u -m
USER u

RUN git clone -q --depth 1 https://github.com/msikma/pokesprite /tmp/original/pokesprite

# Convert all of the pokesprite .pngs -> cowfiles for the terminal
RUN go run /usr/local/src/src/bin/convert/png_convert.go \
        -from /tmp/original/pokesprite/ \
        -to /tmp/cows/ \
        -padding 4 \
        -skip '["resources/", "misc/", "icons/", "items/", "items-outline/"]' \
    && mv -v /tmp/cows/pokemon-gen8 /tmp/cows/gen8 \
    && mv -v /tmp/cows/pokemon-gen7x /tmp/cows/gen7x \
    && cat /tmp/original/pokesprite/data/pokemon.json | jq -c .[] > /tmp/cows/pokemon.json \
    && rm -rf /tmp/original/pokesprite
ADD . .
