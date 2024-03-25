# Building binaries

## On your host OS

_Dependencies_: `Go 1.19`

You'll have to install a golang version that matches the go.mod, and ensure that other package
dependencies are installed (see the dockerfile for the dependencies)

```shell
# Generate binary asset files from the cowfiles
./build/build_assets.sh

# Finally, build the pokesay tool
go build pokesay.go
```

## In docker

_Dependencies:_ `docker`

In order to re/build the binaries from scratch, along with all the cowfile conversion, use the handy
Makefile tasks

```shell
make -C build build/docker build/assets build/release
```

This will produce 4 executable bin files inside the `build/bin` directory, and a heap of binary asset files in `build/assets`.
