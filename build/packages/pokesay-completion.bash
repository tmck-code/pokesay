# pokesay bash completion

_pokesay_completions()
{
    local cur prev names ids cats
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    local opts=(
        --help -h --verbose -v --name -n --id -i --category -c --list-names -l --list-categories -L --width -w --tab-width -t --no-wrap -W --no-tab-spaces -s --fastest -f --no-bubble -B --japanese-name -j --id-info -I --no-category-info -C --info-border -b --unicode-borders -u --flip -F completion
    )

    # Read names and ids from external files
    names="$(</usr/share/pokesay/pokesay-names.txt)"
    ids="$(</usr/share/pokesay/pokesay-ids.txt)"
    cats="big female gen7x gen8 medium regular right shiny small"

    if [[ ${cur} == -* ]] ; then
        COMPREPLY=( $(compgen -W "${opts[*]}" -- ${cur}) )
        return 0
    fi

    if [[ ${prev} == "--category" || ${prev} == "-c" ]]; then
        COMPREPLY=( $(compgen -W "${cats}" -- ${cur}) )
        return 0
    fi

    if [[ ${prev} == "--name" || ${prev} == "-n" ]]; then
        COMPREPLY=( $(compgen -W "${names}" -- ${cur}) )
        return 0
    fi

    if [[ ${prev} == "--id" || ${prev} == "-i" ]]; then
        COMPREPLY=( $(compgen -W "${ids}" -- ${cur}) )
        return 0
    fi
}

complete -F _pokesay_completions pokesay
