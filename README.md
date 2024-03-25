# pokesay

Print pokemon in the CLI! An adaptation of the classic "cowsay"

- [pokesay](#pokesay)
  - [One-line installs](#one-line-installs)
  - [Usage](#usage)
    - [Full Usage](#full-usage)
    - [Examples](#examples)
  - [How it works](#how-it-works)
  - [TODO](#todo)
  - [Similar projects](#similar-projects)

**Other docs**

- [Building binaries](./docs/build.md)
- [Developing/Deploying](./docs/development.md)

---

## Demo

<p align="center">
  <img src="./docs/demo.gif" alt="demo"/>
</p>

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

### Full Usage

> Run pokesay with `-h` or `--help` to see the full usage

```shell
Usage: pokesay [-bCfhjLlsuvW] [-c value] [-n value] [-t value] [-w value] [parameters ...]
 -b, --info-border  draw a border around the info box
 -c, --category=value
                    choose a pokemon from a specific category
 -C, --no-category-info
                    do not print pokemon category information in the info box
 -f, --fastest      run with the fastest possible configuration (--nowrap &
                    --notabspaces)
 -h, --help         display this help message
 -j, --japanese-name
                    print the japanese name in the info box
 -L, --list-categories
                    list all available categories
 -l, --list-names   list all available names
 -n, --name=value   choose a pokemon from a specific name
 -s, --no-tab-spaces
                    do not replace tab characters (fastest)
 -t, --tab-width=value
                    replace any tab characters with N spaces [4]
 -u, --unicode-borders
                    use unicode characters to draw the border around the speech
                    box (and info box if --info-border is enabled)
 -v, --verbose      print verbose output
 -W, --no-wrap      disable text wrapping (fastest)
 -w, --width=value  the max speech bubble width [80]
```

- List all available categories
  ```shell
  pokesay -L
  ```
- List all available names
  ```shell
  pokesay -l
  ```
- Print a message with a random pokemon
  ```shell
  echo 'Hello, world!' | pokesay
  ```
- Print a message with a specific pokemon
  ```shell
  echo 'Hello, world!' | pokesay -n pikachu
  ```
- Print a message with a specific pokemon category
  ```shell
  # big pokemon (i.e. with a large dimensions in the terminal)
  echo 'Hello, world!' | pokesay -c big
  # shiny pokemon
  echo 'Hello, world!' | pokesay -c shiny
  ```
- Print a message with a specific pokemon category and name
  ```shell
  # for shiny charizards
  echo 'Hello, world!' | pokesay -c shiny -n charizard
  ```

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
  - [ ] optionally print ID assigned to each pokemon, support deterministic selection via the same ID
- Longer-term
- In Beta
  - [x] support long and short cli args (e.g. --name/-n)
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

## Similar projects

There are many other projects that bring pokemon to the terminal!
Check them out via the links.

Inspired by the [pokeshell](https://github.com/acxz/pokeshell) project, I've included a comparison table

| project | dependencies | speed | japanese names | size categories | selection by name | selection by category | animated sprites |
|--|--|--|--|--|--|--|--|
| [pokeshell](https://github.com/acxz/pokeshell) | imagemagick, chafa | ? | ? | ? | ? | ? | ? |
| [pokemon-icat](https://github.com/aflaag/pokemon-icat) | python | ? | ? | ? | ? | ? | ? |
| [pokemon-colorscripts](https://gitlab.com/phoneybadger/pokemon-colorscripts) | python3 | ? | ? | ? | ? | ? | ? |
| [pokemonsay-newgenerations](https://github.com/HRKings/pokemonsay-newgenerations) | cowsay (perl) | ? | ? | ? | ? | ? | ? |
| [pokeTerm](https://github.com/31marcosalsa/pokeTerm) | python3 | ? | ? | ? | ? | ? | ? |
| [krabby](https://github.com/yannjor/krabby) | rust, cargo | ? | ❌ | ❌ | ✓ | ❌ | ❌ |
| [pokemonsay](https://github.com/dfrankland/pokemonsay) | npm | ? | ? | ? | ? | ? | ? |