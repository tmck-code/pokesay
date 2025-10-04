function __pokesay_complete
    set -l opts --help -h --verbose -v --name -n --id -i --category -c --list-names -l --list-categories -L --width -w --tab-width -t --no-wrap -W --no-tab-spaces -s --fastest -f --no-bubble -B --japanese-name -j --id-info -I --no-category-info -C --info-border -b --unicode-borders -u --flip -F completion
    set -l names (cat /usr/share/pokesay/pokesay-names.txt)
    set -l ids (cat /usr/share/pokesay/pokesay-ids.txt)
    set -l cats big female gen7x gen8 medium regular right shiny small

    for opt in $opts
        complete -c pokesay -l (string replace -- -- '' $opt) -d "$opt option"
    end
    complete -c pokesay -l name -d "Choose a Pokémon name" -a "$names" -r
    complete -c pokesay -l id -d "Choose a Pokémon ID" -a "$ids" -r
    complete -c pokesay -l category -d "Choose a category" -a "$cats" -r
    complete -c pokesay -s n -d "Choose a Pokémon name" -a "$names" -r
    complete -c pokesay -s i -d "Choose a Pokémon ID" -a "$ids" -r
    complete -c pokesay -s c -d "Choose a category" -a "$cats" -r
end

__pokesay_complete
