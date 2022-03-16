# pokesay-go
A (much) faster go version of tmck-code/pokesay

## One-line installs

- OSX / darwin
    ```shell
    bash -c "$(curl https://raw.githubusercontent.com/tmck-code/pokesay-go/master/scripts/install.sh)" bash darwin amd64
    ```
- Linux x64
    ```shell
    bash -c "$(curl https://raw.githubusercontent.com/tmck-code/pokesay-go/master/scripts/install.sh)" bash linux amd64
    ```
- Android ARM64 (Termux)
    ```shell
    bash -c "$(curl https://raw.githubusercontent.com/tmck-code/pokesay-go/master/scripts/install.sh)" bash android arm64
    ```

---

## How it works

This project extends on the original `fortune | cowsay`, a simple command combo that can be added to your .bashrc to give you a random message spoken by a cow every time you open a new shell.

```
 ☯ ~ fortune | cowsay
 ______________________________________
/ Hollywood is where if you don't have \
| happiness you send out for it.       |
|                                      |
\ -- Rex Reed                          /
 --------------------------------------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||
```

Thee pokemon sprites used here are sourced from the awesome repo [msikma/pokesprite](https://github.com/msikma/pokesprite)

![sprits](https://github.com/msikma/pokesprite/raw/master/resources/images/banner_gen8_2x.png)

All of these sprites are converted into a form that can be rendered in a terminal (unicode characters and colour control sequences) by the `img2xterm` tool, found at [rossy/img2xterm](https://github.com/rossy/img2xterm)

The last pre-compile step is to use `encoding/gob` and `go:embed` to generate a go source code file that encodes all of the converted unicode sprites as gzipped text.

Finally, this is built with the main CLI logic in `pokesay.go` into an single executable that can be easily popped into a directory in the user's $PATH

If all you are after is installing the program to use, then there are no dependencies required! Navigate to the Releases and download the latest binary.

## Building binaries

### In docker

_Dependencies:_ `docker`

In order to re/build the binaries from scratch, along with all the cowfile conversion, use the handy Makefile tasks

```shell
cd build && make build/docker build/release
```

This will produce 4 executable bin files inside the build/ directory

### On your host OS

You'll have to install a golang version that matches the go.mod, and ensure that other package dependencies are installed (see the dockerfile for the dependencies)

```
# Build the pokedex build tool
go build src/pokedex.go
# Extract the cowfiles into a directory
tar xzf build/cows.tar.gz -C build/

# Generate a encoding/gob data file from the cowfiles
./pokedex -from build/cows -to build/cows.gob

# Finally, build the pokesay tool (this builds and uses the build/cows.gob file automatically)
go build pokesay.go
```
