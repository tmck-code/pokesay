# pokesay

Print pokemon in the CLI! An adaptation of the classic "cowsay"

- [pokesay](#pokesay)
  - [One-line installs](#one-line-installs)
  - [Usage](#usage)
  - [How it works](#how-it-works)
  - [Building binaries](#building-binaries)
    - [On your host OS](#on-your-host-os)
    - [In docker](#in-docker)
  - [Developing](#developing)
    - [Testing](#testing)
  - [TODO](#todo)

---

![pokesay demo](https://user-images.githubusercontent.com/9894426/212685620-762e87be-248c-4276-bae6-282c2991db68.png)

---

## One-line installs

_(These commands can also be used to update your existing pokesay)_

- OSX / darwin
    ```shell
    bash -c "$(curl https://raw.githubusercontent.com/tmck-code/pokesay/master/build/scripts/install.sh)" bash darwin amd64
    ```
- OSX / darwin (M1)
    ```shell
    bash -c "$(curl https://raw.githubusercontent.com/tmck-code/pokesay/master/build/scripts/install.sh)" bash darwin arm64
    ```
- Linux x64
    ```shell
    bash -c "$(curl https://raw.githubusercontent.com/tmck-code/pokesay/master/build/scripts/install.sh)" bash linux amd64
    ```
- Android ARM64 (Termux)
    ```shell
    bash -c "$(curl https://raw.githubusercontent.com/tmck-code/pokesay/master/build/scripts/install.sh)" bash android arm64
    ```
- Windows x64 (.exe)
    ```shell
    bash -c "$(curl https://raw.githubusercontent.com/tmck-code/pokesay/master/build/scripts/install.sh)" bash windows amd64
    ```

---

## Usage

Just pipe some text! e.g.

```shell
echo yolo | pokesay
```

To see it every time you open a terminal, add it to your `.bashrc` file!   
_(This requires that you have `fortune` installed)_

```shell
echo 'fortune | pokesay' >> $HOME/.bashrc
```

> _Note: The pokesay tool is intended to only be used with piped text input from STDIN, entering text by typing (or other methods) might not work as expected!_

---

## How it works

This project extends on the original `fortune | cowsay`, a simple command combo that can be added to
your .bashrc to give you a random message spoken by a cow every time you open a new shell.

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

As a personal project, this has been lovingly over-engineered with a focus on the lowest latency possible, so that it doesn't slow down your terminal experience.

1. These pokemon sprites used here are sourced from the awesome repo
[msikma/pokesprite](https://github.com/msikma/pokesprite)

    ![sprits](https://github.com/msikma/pokesprite/raw/master/resources/images/banner_gen8_2x.png)

2. All of these sprites are converted into a form that can be rendered in a terminal (unicode
characters and colour control sequences) by the `img2xterm` tool, found at
[rossy/img2xterm](https://github.com/rossy/img2xterm)

3. Use some go tools (`encoding/gob` and `go:embed`) to generate a go source code file
that encodes all of the converted unicode sprites as gzipped text and some search-optimised data structures.

4. Finally, this is built with the main CLI logic in `pokesay.go` into an single executable that can be
easily popped into a directory in the user's `$PATH`

If all you are after is installing the program to use, then there are no dependencies required!
Navigate to the Releases and download the latest binary.

---

## Building binaries

### On your host OS

_Dependencies_: `Go 1.19`

You'll have to install a golang version that matches the go.mod, and ensure that other package
dependencies are installed (see the dockerfile for the dependencies)

```shell
# Generate binary asset files from the cowfiles
./build/build_assets.sh

# Finally, build the pokesay tool
go build pokesay.go
```

### In docker

_Dependencies:_ `docker`

In order to re/build the binaries from scratch, along with all the cowfile conversion, use the handy
Makefile tasks

```shell
make -C build build/docker build/assets build/release
```

This will produce 4 executable bin files inside the `build/bin` directory, and a heap of binary asset files in `build/assets`.

---

## Developing

### Testing

Run tests with

```shell
go test -v ./test/

# or, with debug information printed
DEBUG=test go test -v ./test/
```

---

## TODO

- Short-term
  - [ ] add short versions of command line args
- Longer-term
- Completed
  - [x] Make the category struct faster to load - currently takes up to 80% of the execution time
  - [x] Store metadata and names in a more storage-efficient manner
  - [x] Import japanese names from data/pokemon.json
  - [x] Fix bad whitespace stripping when building assets
  - [x] List all names
  - [x] Make data structure to hold categories, names and pokemon
  - [x] Increase speed
  - [x] Improve categories to be more specific than shiny/regular
