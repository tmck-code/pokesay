#compdef pokesay

_pokesay() {
  local -a opts names ids cats
  opts=(
    --help -h --verbose -v --name -n --id -i --category -c --list-names -l --list-categories -L --width -w --tab-width -t --no-wrap -W --no-tab-spaces -s --fastest -f --no-bubble -B --japanese-name -j --id-info -I --no-category-info -C --info-border -b --unicode-borders -u --flip -F completion
  )
  names=($(</usr/share/pokesay/pokesay-names.txt))
  ids=($(</usr/share/pokesay/pokesay-ids.txt))
  cats=(big female gen7x gen8 medium regular right shiny small)

  _arguments \
    "${opts[@]/#/-}" \
    '--name[Choose a Pokémon name]:name:(( ${names[*]} ))' \
    '--id[Choose a Pokémon ID]:id:(( ${ids[*]} ))' \
    '--category[Choose a category]:category:(( ${cats[*]} ))' \
    '--list-names[List all available names]' \
    '--list-categories[List all available categories]' \
    '--width[Set max speech bubble width]' \
    '--tab-width[Set tab width]' \
    '--no-wrap[Disable text wrapping]' \
    '--no-tab-spaces[Do not replace tab characters]' \
    '--fastest[Fastest configuration]' \
    '--no-bubble[Do not draw speech bubble]' \
    '--japanese-name[Print Japanese name]' \
    '--id-info[Print Pokémon ID]' \
    '--no-category-info[Do not print category info]' \
    '--info-border[Draw info box border]' \
    '--unicode-borders[Use unicode borders]' \
    '--flip[Flip Pokémon]' \
    '-n[Choose a Pokémon name]:name:(( ${names[*]} ))' \
    '-i[Choose a Pokémon ID]:id:(( ${ids[*]} ))' \
    '-c[Choose a category]:category:(( ${cats[*]} ))' \
    '-l[List all available names]' \
    '-L[List all available categories]' \
    '-w[Set max speech bubble width]' \
    '-t[Set tab width]' \
    '-W[Disable text wrapping]' \
    '-s[Do not replace tab characters]' \
    '-f[Fastest configuration]' \
    '-B[Do not draw speech bubble]' \
    '-j[Print Japanese name]' \
    '-I[Print Pokémon ID]' \
    '-C[Do not print category info]' \
    '-b[Draw info box border]' \
    '-u[Use unicode borders]' \
    '-F[Flip Pokémon]'
}

compdef _pokesay pokesay
