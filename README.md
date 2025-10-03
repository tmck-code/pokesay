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

|   |   |   |
|---|---|---|
| <img width="300" alt="p00" src="https://github.com/user-attachments/assets/202d7142-8ced-4bd9-afef-53024e215ea3" /> | <img width="300" alt="p01" src="https://github.com/user-attachments/assets/9916020d-bb32-45b6-b701-b1a9e0c010d3" /> | <img width="300" alt="p02" src="https://github.com/user-attachments/assets/a060690f-a12a-46e4-b16f-a2b88242da78" /> |
| <img width="300" alt="p03" src="https://github.com/user-attachments/assets/904a9e99-f6f9-4ab7-a9fc-fae7b1a60780" /> | <img width="300" alt="p04" src="https://github.com/user-attachments/assets/603bb240-ace8-433d-b85a-e9e78db081d0" /> | <img width="300" alt="p05" src="https://github.com/user-attachments/assets/e45bd570-bad3-4bb6-92e9-06aa318f3cae" /> |
| <img width="300" alt="p06" src="https://github.com/user-attachments/assets/b5603034-3d36-4257-8948-a4258698fc5d" /> | <img width="300" alt="p07" src="https://github.com/user-attachments/assets/30e425f7-7fbb-44a6-958a-2046dc7c4bb3" /> | <img width="300" alt="p08" src="https://github.com/user-attachments/assets/d88d4ed0-a47a-458f-8eca-b397e0552569" /> |
| <img width="300" alt="p09" src="https://github.com/user-attachments/assets/9dfe7f26-46e4-4eb0-b476-02571a3a64af" /> | <img width="300" alt="p10" src="https://github.com/user-attachments/assets/d928cf62-0a48-43b9-8640-c3f9255fc7b3" /> | <img width="300" alt="p11" src="https://github.com/user-attachments/assets/b1b1c7dc-cd93-4993-86e8-38d8bd69ab5f" /> |
| <img width="300" alt="p12" src="https://github.com/user-attachments/assets/153c804f-3251-471d-9887-9e92b1955a2c" /> | <img width="300" alt="p13" src="https://github.com/user-attachments/assets/b42b009c-5b80-4816-99a7-db0d4248a6b0" /> | <img width="300" alt="p14" src="https://github.com/user-attachments/assets/41005ff1-8fb0-43cb-99be-8e662818479f" /> |
| <img width="300" alt="p15" src="https://github.com/user-attachments/assets/7e537507-8b9b-4447-8061-a700d3f83412" /> |  |  |

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
