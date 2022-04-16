FROM golang:latest

WORKDIR /usr/local/src

RUN BUILD_DEPS="parallel make gcc libmagickwand-dev ncurses-dev fortune-mod fortunes" \
    && apt-get update \
    && apt-get install -y --no-install-recommends ${BUILD_DEPS} \
    && git clone --depth 1 https://github.com/denilsonsa/img2xterm \
      && (cd img2xterm && make && make install) \
    && rm -rf imgxterm/ \
    && apt-get purge -y --auto-remove ${BUILD_DEPS} \
    && rm -rf /var/lib/apt/lists/*

ADD go.mod go.sum src/ /usr/local/src/

# Convert all of the pokesprite .pngs -> cowfiles for the terminal
RUN git clone --depth 1 https://github.com/msikma/pokesprite /tmp/original/pokesprite \
    && go mod tidy \
    && go run /usr/local/src/src/bin/convert/png_convert.go \
        -from /tmp/original/pokesprite/ \
        -to /tmp/cows/ \
        -padding 4 \
        -skip '["resources/", "misc/", "icons/", "items/", "items-outline/"]' \
    && rm -rf /tmp/original

ADD . .
