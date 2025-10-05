#compdef pokesay

_pokesay() {
  local -a opts names ids cats
  # Option descriptions from README Full Usage
  # e.g. 
  # '-T+[Defines the window title \[default\: Alacritty\]]:TITLE:_default' \
  # '--title=[Defines the window title \[default\: Alacritty\]]:TITLE:_default' \
  opts=(
    '-h[Display this help message]:HELP'
    '--help[Display this help message]:HELP'
    '-v[Print verbose output]:VERBOSE'
    '--verbose[Print verbose output]:VERBOSE'
    '-n=[Choose a Pokémon from a specific name]:NAME:(${names[*]})'
    '--name=[Choose a Pokémon from a specific name]:NAME:(${names[*]})'
    '-i=[Choose a Pokémon from a specific ID]:ID:(${ids[*]})'
    '--id=[Choose a Pokémon from a specific ID]:ID:(${ids[*]})'
    '-c=[Choose a Pokémon from a specific category]:CATEGORY:(${cats[*]})'
    '--category=[Choose a Pokémon from a specific category]:CATEGORY:(${cats[*]})'
    '-l[List all available names]:LIST_NAMES'
    '--list-names[List all available names]:LIST_NAMES'
    '-L[List all available categories]:LIST_CATEGORIES'
    '--list-categories[List all available categories]:LIST_CATEGORIES'
    '-w=[Set max speech bubble width \[default\: 80\]]:WIDTH'
    '--width=[Set max speech bubble width \[default\: 80\]]:WIDTH'
    '-t=[Replace tab characters with N spaces \[default\: 4\]]:TAB-WIDTH'
    '--tab-width=[Replace tab characters with N spaces \[default\: 4\]]:TAB-WIDTH'
    '-W[Disable text wrapping (fastest)]:DISABLE_WRAP'
    '--no-wrap[Disable text wrapping (fastest)]:DISABLE_WRAP'
    '-s[Do not replace tab characters (fastest)]:DISABLE_TAB_SPACES'
    '--no-tab-spaces[Do not replace tab characters (fastest)]:DISABLE_TAB_SPACES'
    '-f[Run with the fastest possible configuration (--nowrap & --notabspaces)]:FASTEST'
    '--fastest[Run with the fastest possible configuration (--nowrap & --notabspaces)]:FASTEST'
    '-B[Do not draw the speech bubble]:DISABLE_BUBBLE'
    '--no-bubble[Do not draw the speech bubble]:DISABLE_BUBBLE'
    '-j[Print the Japanese name in the info box]:JAPANESE_NAME'
    '--japanese-name[Print the Japanese name in the info box]:JAPANESE_NAME'
    '-I[Print the Pokémon ID in the info box]:POKEMON_ID'
    '--id-info[Print the Pokémon ID in the info box]:POKEMON_ID'
    '-C[Do not print category info in the info box]:DISABLE_CATEGORY_INFO'
    '--no-category-info[Do not print category info in the info box]:DISABLE_CATEGORY_INFO'
    '-b[Draw a border around the info box]:INFO_BORDER'
    '--info-border[Draw a border around the info box]:INFO_BORDER'
    '-u[Use unicode characters to draw the border]:UNICODE_BORDERS'
    '--unicode-borders[Use unicode characters to draw the border]:UNICODE_BORDERS'
    '-F[Flip the Pokémon horizontally (face right instead of left)]:FLIP'
    '--flip[Flip the Pokémon horizontally (face right instead of left)]:FLIP'
  )
  names=($(</usr/share/pokesay/pokesay-names.txt))
  ids=($(</usr/share/pokesay/pokesay-ids.txt))
  cats=(big female gen7x gen8 medium regular right shiny small)

  _arguments \
    ${opts[@]}
}

compdef _pokesay pokesay
