FROM golang:latest

WORKDIR /usr/local/src

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        parallel make gcc libmagickwand-dev ncurses-dev fortune-mod fortunes

RUN git clone --depth 1 https://github.com/denilsonsa/img2xterm \
    && (cd img2xterm && make && make install)

WORKDIR /tmp/original

RUN git clone --depth 1 https://github.com/msikma/pokesprite

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
        -skip '["resources/", "misc/", "icons/", "items/", "items-outline/"]'

ADD . .