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

> _Note: The pokesay tool is intended to only be used with piped text input from STDIN, entering text by typing (or other methods) might not work as expected!_

### Examples


<table>
<tr>
  <td valign="top">Default output<br><img width="350" alt="p_00_fortune" src="https://github.com/user-attachments/assets/019a9055-0a71-472e-84df-791ff8803bb4" /></td>
  <td valign="top">Set bubble width <b><code>-w</code></b><br><img width="350" alt="p_01_set_width" src="https://github.com/user-attachments/assets/c458d24d-ea0c-4d36-abd3-1a329e81dfde" /></td>
  <td valign="top">Unicode borders <b><code>-u</code></b><br><img width="350" alt="p_02_unicode" src="https://github.com/user-attachments/assets/9620db09-58e3-4c38-b1a4-99aab6b6e85e" /></td>
</tr>
<tr>
  <td valign="top">Info border <b><code>-b</code></b><br><img width="350" alt="p_03_info_border" src="https://github.com/user-attachments/assets/fe857d59-4400-4251-824b-171d5f10ac52" /></td>
  <td valign="top">Choose by name <b><code>-n</code></b><br><img width="350" alt="p_04_name" src="https://github.com/user-attachments/assets/88e5fc35-28ac-4b48-a04d-fb1f14446d9d" /></td>
  <td valign="top">No category info <b><code>-C</code></b><br><img width="350" alt="p_05_no_category" src="https://github.com/user-attachments/assets/dc9217b8-77a8-4add-a42c-c734d488d1ec" /></td>
</tr>
<tr>
  <td valign="top">Show Japanese name <b><code>-j</code></b><br><img width="350" alt="p_06_japanese" src="https://github.com/user-attachments/assets/34ebce1c-4448-4560-aa78-704cbabaef1e" /></td>
  <td valign="top">Show ID <b><code>-I</code></b><br><img width="350" alt="p_07_id" src="https://github.com/user-attachments/assets/cd2ba951-89db-4764-9dbd-553105d54707" /></td>
  <td valign="top">Small size <b><code>-c small</code></b><br><img width="350" alt="p_08_small_size" src="https://github.com/user-attachments/assets/644d9574-2c67-4e2c-a586-b4b0aa95211c" /></td>
</tr>
<tr>
  <td valign="top">Medium size <b><code>-c medium</code></b><br><img width="350" alt="p_09_medium_size" src="https://github.com/user-attachments/assets/9979c100-76a2-479c-a5fc-cea858906b25" /></td>
  <td valign="top">Big size <b><code>-c big</code></b><br><img width="350" alt="p_10_big_size" src="https://github.com/user-attachments/assets/a65c822c-fcd8-4a1a-bdc0-bb81ff9e164e" /></td>
  <td valign="top">Shiny <b><code>-c shiny</code></b><br><img width="350" alt="p_10_shiny" src="https://github.com/user-attachments/assets/1a0022fb-921f-46b8-911c-8db7a8e8769b" /></td>
</tr>
<tr>
  <td valign="top">Size and name <b><code>-c small -n ...</code></b><br><img width="350" alt="p_11_size_and_name" src="https://github.com/user-attachments/assets/e0dc9d09-b439-4dc5-81b0-ea5419308c47" /></td>
  <td valign="top">Size and name <b><code>-c big -n ...</code></b><br><img width="350" alt="p_12_size_and_name_2" src="https://github.com/user-attachments/assets/339a2832-9488-49bf-89ef-abf0f376e336" /></td>
  <td valign="top"></td>
</tr>
<tr>
  <td valign="top">Select by ID <b><code>-i</code></b><br><img width="350" alt="p_13_select_by_id" src="https://github.com/user-attachments/assets/73aa44c5-bec5-41b0-98b7-717dd29b395f" /></td>
  <td valign="top">Flip <b><code>-F</code></b><br><img width="350" alt="p_14_flip" src="https://github.com/user-attachments/assets/f610c99d-6ad4-4eaa-82ff-0afbe1a7d2a3" /></td>
  <td valign="top">Figlet + lolcat<br><img width="350" alt="p_15_figlet_lolcat" src="https://github.com/user-attachments/assets/1e88880b-adb2-4cc2-9ee9-29c0acf26f57" /></td>
</tr>
</table>


To see it every time you open a terminal, add it to your `.bashrc` file!   
_(This requires that you have `fortune` installed)_

```shell
echo 'fortune | pokesay' >> $HOME/.bashrc
```

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
