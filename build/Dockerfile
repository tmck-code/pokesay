FROM golang:latest

WORKDIR /usr/local/src

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        parallel make gcc libmagickwand-dev ncurses-dev fortune-mod fortunes

RUN git clone --depth 1 https://github.com/denilsonsa/img2xterm \
    && (cd img2xterm && make && make install)

ENV DOCKER_BUILD_DIR /tmp/original
ENV DOCKER_OUTPUT_DIR /tmp/cows

RUN mkdir -p /tmp/original /tmp/cows
WORKDIR /tmp/original

RUN git clone --depth 1 https://github.com/msikma/pokesprite

WORKDIR /usr/local/src/
ADD src/ /usr/local/src/src/
ADD go.* /usr/local/src/
RUN go mod tidy \
    && go get -v github.com/mitchellh/go-wordwrap

ADD . .
RUN mkdir -p build/assets/metadata build/assets/cows

# Convert all of the pokesprite .pngs -> cowfiles for the terminal
RUN go run /usr/local/src/src/bin/convert/png_convert.go \
        -from /tmp/original/pokesprite/ \
        -to /tmp/cows/ \
        -skip '["resources/", "misc/", "icons/", "items/", "items-outline/"]'

# generate embedded bin files for category/metadata/the actual pokemon
RUN go run /usr/local/src/src/bin/pokedex/pokedex.go \
    -from /tmp/cows/ \
    -to build/assets/ \
    -toCategoryFpath pokedex.gob \
    -toDataSubDir cows/ \
    -toMetadataSubDir metadata/ \
    -toTotalFname total.txt
