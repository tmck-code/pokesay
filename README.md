# pokesay

Print pokemon in the CLI! An adaptation of the classic 'cowsay'

<img width="1065" alt="464197914-4c145f23-8837-41df-835c-aeaa49afd13d 2" src="https://github.com/user-attachments/assets/57d33b92-95cf-4a5b-890c-39ff530d447c" />

---

## Installation

- Via homebrew (MacOS/Linux/Windows)
  ```shell
  brew install tmck-code/tap/pokesay
  ```
- Via the AUR (Arch Linux)
  ```shell
  yay -S pokesay-bin
  ```

Pokesay is a single binary with no dependencies that can be run on `arm64/amd64 OSX`, `amd64 Linux`, `arm64 Android` and `amd64 Windows`.   
For installation without a package manager, see the options below.

<details>
<summary><i>Via the install script (Others)</i></summary>

_(These commands can also be used to update your existing pokesay)_

<table>
<tr>
  <td>OS/arch</td> <td>command</td>
</tr>
<tr>
  <td>OSX / darwin</td>
  <td>

```shell
bash -c "$(curl https://raw.githubusercontent.com/tmck-code/pokesay/master/build/scripts/install.sh)" \
  bash darwin amd64
```

  </td>
</tr>
<tr>
  <td>OSX / darwin (M1)</td>
  <td>

```shell
bash -c "$(curl https://raw.githubusercontent.com/tmck-code/pokesay/master/build/scripts/install.sh)" \
  bash darwin arm64
```

  </td>
</tr>
<tr>
  <td>Linux / x64</td>
  <td>

```shell
bash -c "$(curl https://raw.githubusercontent.com/tmck-code/pokesay/master/build/scripts/install.sh)" \
  bash linux amd64
```

  </td>
</tr>
<tr>
  <td>Android / arm64 (termux)</td>
  <td>

```shell
bash -c "$(curl https://raw.githubusercontent.com/tmck-code/pokesay/master/build/scripts/install.sh)" \
  bash android arm64
```

  </td>
</tr>
<tr>
  <td>Windows / x64 (.exe)</td>
  <td>

```shell
bash -c "$(curl https://raw.githubusercontent.com/tmck-code/pokesay/master/build/scripts/install.sh)" \
  bash windows amd64
```

  </td>
</tr>
</table>

</details>

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

<p align="center">
  <kbd><img src="https://github.com/tmck-code/pokesay/assets/9894426/794ee42d-4bc2-4d0c-bce0-aafa3dca2e78" alt="demo"/></kbd>
</p>

### Full Usage

> Run pokesay with `-h` or `--help` to see the full usage

```shell
Usage: pokesay [-bBCfFhIjLsuvW] [-c value] [-i value] [-l value] [-n value] [-t value] [-w value] [parameters ...]
 -b, --info-border  draw a border around the info box
 -B, --no-bubble    do not draw the speech bubble
 -c, --category=value
                    choose a pokemon from a specific category
 -C, --no-category-info
                    do not print pokemon category information in the info box
 -f, --fastest      run with the fastest possible configuration (--nowrap &
                    --notabspaces)
 -F, --flip         flip the pokemon horizontally (face right instead of left)
 -h, --help         display this help message
 -i, --id=value     choose a pokemon from a specific ID (see `pokesay -l` for
                    IDs)
 -I, --id-info      print the pokemon ID in the info box
 -j, --japanese-name
                    print the japanese name in the info box
 -L, --list-categories
                    list all available categories
 -l, --list-names[=value]
                    list all available names
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

### Examples

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
- Print a specific pokemon by its ID
  ```shell
  # green mewtwo is ID `0491.1719`
  echo 'Hello, world!' | pokesay -i 0491.1719
  ```

---

<img width="1000" height="1008" alt="p18" src="https://github.com/user-attachments/assets/5a1b8e27-15f8-40df-9f3e-6117dd5ebda5" />
<img width="550" height="670" alt="p17" src="https://github.com/user-attachments/assets/146d0fa1-8d82-4350-9a12-c3fe9840b78a" />
<img width="550" height="670" alt="p16" src="https://github.com/user-attachments/assets/7e080bb1-b60d-421d-8f9c-70196f04c8dd" />
<img width="520" height="592" alt="p15" src="https://github.com/user-attachments/assets/6b8dd048-30f0-4dd4-b4fb-8013062dbded" />
<img width="500" height="592" alt="p14" src="https://github.com/user-attachments/assets/0762fa3a-5bd7-4ef8-9b53-bae1c0e7ba2e" />
<img width="500" height="592" alt="p13" src="https://github.com/user-attachments/assets/1b3fe721-8561-42e4-9d09-008c602cbe00" />
<img width="500" height="540" alt="p12" src="https://github.com/user-attachments/assets/b51a5b0a-2c64-4a6f-a55a-5fb207ddfe7d" />
<img width="500" height="592" alt="p11" src="https://github.com/user-attachments/assets/8a7252d2-a1a5-4005-883a-685bfa285dbb" />
<img width="500" height="592" alt="p10" src="https://github.com/user-attachments/assets/3c89e6ab-3910-45fc-a4d0-e9bb5c4a3c9f" />
<img width="600" height="800" alt="p09" src="https://github.com/user-attachments/assets/c2d3014e-5192-475d-810d-c034cf394756" />
<img width="550" height="592" alt="p08" src="https://github.com/user-attachments/assets/d536097b-1243-4f90-841a-b5e10496b7b9" />
<img width="600" height="800" alt="p07" src="https://github.com/user-attachments/assets/084fc382-fac9-4d18-a7b3-941f8525f99d" />
<img width="500" height="670" alt="p06" src="https://github.com/user-attachments/assets/e0bb0c74-6fae-4c84-a436-389d9c0ac27d" />
<img width="450" height="540" alt="p05" src="https://github.com/user-attachments/assets/bdd9b63c-3b3d-422a-819a-c0b1529258f7" />
<img width="450" height="566" alt="p04" src="https://github.com/user-attachments/assets/d1c69b37-2aae-4879-b957-0b79d694d868" />
<img width="450" height="592" alt="p03" src="https://github.com/user-attachments/assets/653002a8-5b8f-44c0-9387-5777938022e7" />
<img width="450" height="514" alt="p02" src="https://github.com/user-attachments/assets/f929038c-4a45-45f2-a1e8-7c3aee9a5ecf" />
<img width="450" height="514" alt="p01" src="https://github.com/user-attachments/assets/4d94d2d0-6b88-4fa6-9e59-50341cfbb369" />
<img width="900" height="566" alt="p00" src="https://github.com/user-attachments/assets/e8c4dfda-a007-4969-8b93-4a75ac5fc21d" />


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

## Similar projects

There are many other projects that bring pokemon to the terminal!
Check them out via the links.

Inspired by the [pokeshell](https://github.com/acxz/pokeshell) project, I've included a comparison table

| project | dependencies | speed | japanese names | size categories | selection by name | selection by category | animated sprites |
|--|--|--|--|--|--|--|--|
| [tmck-code/pokesay](https://github.com/tmck-code/pokesay) | none 🎉 | ~2.5ms | ✅ | ✅ | ✅ | ✅ | ❌ |
| [pokeshell](https://github.com/acxz/pokeshell) | imagemagick, chafa | ? | ? | ? | ? | ? | ? |
| [pokemon-icat](https://github.com/aflaag/pokemon-icat) | python | ? | ? | ? | ? | ? | ? |
| [pokemon-colorscripts](https://gitlab.com/phoneybadger/pokemon-colorscripts) | python3 | ? | ? | ? | ? | ? | ? |
| [pokemonsay-newgenerations](https://github.com/HRKings/pokemonsay-newgenerations) | cowsay (perl) | ? | ? | ? | ? | ? | ? |
| [31marcosalsa/pokeTerm](https://github.com/31marcosalsa/pokeTerm) | python, imagemagick, img2xterm | ~62ms | ❌ | ❌ | ❌ | ❌ | ❌ |
| [krabby](https://github.com/yannjor/krabby) | rust, cargo | ? | ❌ | ❌ | ✓ | ❌ | ❌ |
| [pokemonsay](https://github.com/dfrankland/pokemonsay) | npm | ? | ? | ? | ? | ? | ? |
| [possatti/pokemonsay](https://github.com/possatti/pokemonsay) | none 🎉 | ~26.3ms | ❌ | ❌ | ✅ | ❌ | ❌ | 

---

## TODO

- **In progress**
- **Short-term**
  - [ ] requesting mew returns mewtwo also
  - [ ] create "vertical" friendly display mode, place the Pokemon standing beside the text box, on the left or right
- **Longer-term**
  - [ ] make the process async.
    - (Currently the searching/pokemon fetching is done _before_ any printing begins. There's an opportunity to start printing the speech bubble while also fetching the pokemon to print below it)
    - [ ] implement native lolcat/rainbow HR/colour
- **In Beta**
  - [x] optionally print ID assigned to each pokemon, support deterministic selection via the same ID
- **Completed**
  - [x] add option to flip Pokemon to face right or left, remove all "right" facing cowfiles
  - [x] create Debian package
  - [x] create Arch package
  - [x] support long and short cli args (e.g. --name/-n)
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
