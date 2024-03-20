# pokesay

Print pokemon in the CLI! An adaptation of the classic "cowsay"

- [One-line installs](#one-line-installs)
- [Usage](#usage)
- [How it works](#how-it-works)
- [TODO](#todo)
- [Other docs](#other-docs)

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

## TODO

- Short-term
  - [ ] support long and short cli args (e.g. --name/-n)
  - [ ] optionally print ID assigned to each pokemon, support deterministic selection via the same ID
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
  - [x] Filter by both name and category

## Other docs

- [Building binaries](./docs/build.md)
- [Developing/Deploying](./docs/development.md)
