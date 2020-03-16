FROM golang:latest

WORKDIR /usr/local/src

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        parallel make gcc libmagickwand-dev ncurses-dev

ENV PATH "/usr/lib/x86_64-linux-gnu/ImageMagick-6.9.10/bin-q16:/usr/include/ImageMagick-6/:$PATH"

RUN git clone --depth 1 git://github.com/msikma/pokesprite && \
    cp -R pokesprite/icons/ . && \
    rm -rf pokesprite && \
    git clone --depth 1 git://github.com/rossy/img2xterm && \
    (cd img2xterm && make && make install)

WORKDIR /usr/local/src
RUN git clone --depth 1 git://github.com/msmith491/go-cowsay

ADD ops/build_cows.sh ops/build_bin.sh ./
