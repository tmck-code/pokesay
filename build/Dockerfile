FROM golang:1.24

WORKDIR /usr/local/src

ENV DEBIAN_FRONTEND=noninteractive
ENV BUILD_DEPS="make gcc libmagickwand-dev ncurses-dev"
RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y --no-install-recommends \
        $BUILD_DEPS jq imagemagick

RUN git clone -q --depth 1 https://github.com/denilsonsa/img2xterm \
    && (cd img2xterm && make && make install) \
    && rm -rf img2xterm \
    && apt-get purge -y $BUILD_DEPS \
    && apt-get autoremove --purge -y \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /tmp/* /var/tmp/*

RUN useradd u -m 
ADD --chown=u:u src/ /home/u/pokesay/src/
ADD --chown=u:u go.* /home/u/pokesay/
WORKDIR /home/u/pokesay/
USER u
ENV GOPATH=/home/u/go
ENV PATH=$GOPATH/bin:$PATH

RUN go mod tidy \
    && go get -v github.com/mitchellh/go-wordwrap \
    && go install gotest.tools/gotestsum@latest

RUN git clone -q --depth 1 https://github.com/msikma/pokesprite /tmp/original/pokesprite

# Convert all of the pokesprite .pngs -> cowfiles for the terminal
RUN go run /home/u/pokesay/src/bin/convert/png_convert.go \
        -from /tmp/original/pokesprite/ \
        -to /tmp/cows/ \
        -padding 4 \
        -skip '["resources/", "misc/", "icons/", "items/", "items-outline/"]' \
    && mv -v /tmp/cows/pokemon-gen8 /tmp/cows/gen8 \
    && mv -v /tmp/cows/pokemon-gen7x /tmp/cows/gen7x \
    && cat /tmp/original/pokesprite/data/pokemon.json | jq -c .[] > /tmp/cows/pokemon.json \
    && rm -rf /tmp/original/pokesprite

ADD --chown=u:u . .
