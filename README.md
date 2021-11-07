# pokesay-go
A (much) faster go version of tmck-code/pokesay

## One-line installs

### OSX / darwin
```shell
bash -c "$(curl https://raw.githubusercontent.com/tmck-code/pokesay-go/master/install.sh)" bash darwin amd64
```

### Linux x64
```shell
bash -c "$(curl https://raw.githubusercontent.com/tmck-code/pokesay-go/master/install.sh)" bash linux amd64
```

### Android ARM64 (Termux)
```shell
bash -c "$(curl https://raw.githubusercontent.com/tmck-code/pokesay-go/master/install.sh)" bash android arm64
```

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

The last pre-compile step is to use `go-bindata` tool ([go-bindata/go-bindata](https://github.com/go-bindata/go-bindata)) to generate a go source code file that encodes all of the converted unicode sprites as binary text.

Finally, this is built with the main CLI logic in `pokesay.go` into an single executable that can be easily popped into a directory in the user's $PATH

If all you are after is installing the program to use, then there are no dependencies required! Navigate to the Releases and download the latest binary.

## Building binaries

_Dependencies:_ `docker`

In order to re/build the binaries from scratch, along with all the cowfile conversion, use the command that matches your target GOOS/GOARCH 

```shell
TARGET_GOOS=darwin TARGET_GOARCH=amd64 make build/bin
TARGET_GOOS=linux TARGET_GOARCH=amd64 make build/bin
TARGET_GOOS=windows TARGET_GOARCH=amd64 make build/bin
TARGET_GOOS=android TARGET_GOARCH=arm64 make build/bin
```
