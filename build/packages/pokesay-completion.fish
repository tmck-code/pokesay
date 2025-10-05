# pokesay fish completions

function __pokesay_complete
    complete -c pokesay -s b -l info-border        -d "Draw a border around the info box"
    complete -c pokesay -s B -l no-bubble          -d "Do not draw the speech bubble"
    complete -c pokesay -s c -l category           -d "Choose a Pokémon from a specific category" -a "$cats" -r
    complete -c pokesay -s C -l no-category-info   -d "Do not print category info in the info box"
    complete -c pokesay -s f -l fastest            -d "Run with the fastest possible configuration (--nowrap & --notabspaces)"
    complete -c pokesay -s F -l flip               -d "Flip the Pokémon horizontally (face right instead of left)"
    complete -c pokesay -s h -l help               -d "Display this help message"
    complete -c pokesay -s i -l id                 -d "Choose a Pokémon from a specific ID" -a "$ids" -r
    complete -c pokesay -s I -l id-info            -d "Print the Pokémon ID in the info box"
    complete -c pokesay -s j -l japanese-name      -d "Print the Japanese name in the info box"
    complete -c pokesay -s L -l list-categories    -d "List all available categories"
    complete -c pokesay -s l -l list-names         -d "List all available names"
    complete -c pokesay -s n -l name               -d "Choose a Pokémon from a specific name" -a "$names" -r
    complete -c pokesay -s s -l no-tab-spaces      -d "Do not replace tab characters (fastest)"
    complete -c pokesay -s t -l tab-width          -d "Replace tab characters with N spaces [4]"
    complete -c pokesay -s u -l unicode-borders    -d "Use unicode characters to draw the border"
    complete -c pokesay -s v -l verbose            -d "Print verbose output"
    complete -c pokesay -s W -l no-wrap            -d "Disable text wrapping (fastest)"
    complete -c pokesay -s w -l width              -d "Set max speech bubble width [80]"

    set -l names (cat /usr/share/pokesay/pokesay-names.txt)
    set -l ids (cat /usr/share/pokesay/pokesay-ids.txt)
    set -l cats big female gen7x gen8 medium regular right shiny small
end

__pokesay_complete
