FROM golang:latest

WORKDIR /usr/local/src

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        parallel make gcc libmagickwand-dev ncurses-dev

RUN git clone --depth 1 git://github.com/msikma/pokesprite && \
    cp -R pokesprite/icons/ . && \
    rm -rf pokesprite && \
    git clone --depth 1 git://github.com/rossy/img2xterm && \
    cd img2xterm && \
    PATH="/usr/lib/x86_64-linux-gnu/ImageMagick-6.9.10/bin-q16:/usr/include/ImageMagick-6/:$PATH" make && \
    PATH="/usr/lib/x86_64-linux-gnu/ImageMagick-6.9.10/bin-q16:/usr/include/ImageMagick-6/:$PATH" make install

WORKDIR /usr/local/src

ADD ops/build.sh .
