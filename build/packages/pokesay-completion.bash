# pokesay bash completion

_pokesay_completions()
{
    local cur prev names ids cats
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    local opts=(
        -h --help
        -v --verbose
        -n --name
        -i --id
        -c --category
        -l --list-names
        -L --list-categories
        -w --width
        -t --tab-width
        -W --no-wrap
        -s --no-tab-spaces
        -f --fastest
        -B --no-bubble
        -j --japanese-name
        -I --id-info
        -C --no-category-info
        -b --info-border
        -u --unicode-borders
        -F --flip
    )

    names="$(</usr/share/pokesay/pokesay-names.txt)"
    ids="$(</usr/share/pokesay/pokesay-ids.txt)"
    cats="big female gen7x gen8 medium regular right shiny small"

    if [[ ${cur} == -* ]] ; then
        COMPREPLY=( $(compgen -W "${opts[*]}" -- ${cur}) )
        return 0
    elif [[ ${prev} == "--category" || ${prev} == "-c" ]]; then
        COMPREPLY=( $(compgen -W "${cats}" -- ${cur}) )
        return 0
    elif [[ ${prev} == "--name" || ${prev} == "-n" ]]; then
        COMPREPLY=( $(compgen -W "${names}" -- ${cur}) )
        return 0
    elif [[ ${prev} == "--id" || ${prev} == "-i" ]]; then
        COMPREPLY=( $(compgen -W "${ids}" -- ${cur}) )
        return 0
    fi
}

complete -F _pokesay_completions pokesay
