FROM golang:1.19

WORKDIR /usr/local/src

RUN BUILD_DEPS="make gcc"
RUN apt-get update \
    && apt-get install -y --no-install-recommends "${BUILD_DEPS}" libmagickwand-dev ncurses-dev jq

RUN git clone --depth 1 https://github.com/denilsonsa/img2xterm \
    && (cd img2xterm && make && make install) \
    && rm -rf img2xterm \
    && apt purge -y "${BUILD_DEPS}" \
    && apt-get purge -y --auto-remove \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /tmp/* /var/tmp/*

RUN git clone --depth 1 https://github.com/msikma/pokesprite /tmp/original/pokesprite

WORKDIR /usr/local/src/
ADD src/ /usr/local/src/src/
ADD go.* /usr/local/src/
RUN go mod tidy \
    && go get -v github.com/mitchellh/go-wordwrap

# Convert all of the pokesprite .pngs -> cowfiles for the terminal
RUN go run /usr/local/src/src/bin/convert/png_convert.go \
        -from /tmp/original/pokesprite/ \
        -to /tmp/cows/ \
        -padding 4 \
        -skip '["resources/", "misc/", "icons/", "items/", "items-outline/"]' \
    && cat /tmp/original/pokesprite/data/pokemon.json | jq -c .[] > /tmp/cows/pokemon.json \
    && rm -rf /tmp/original/pokesprite

ADD . .
