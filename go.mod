module github.com/tmck-code/pokesay-go

go 1.17

require (
	github.com/mitchellh/go-wordwrap v1.0.1
	github.com/schollz/progressbar/v3 v3.8.3
)

require (
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	golang.org/x/crypto v0.0.0-20211202192323-5770296d904e // indirect
	golang.org/x/sys v0.0.0-20211205182925-97ca703d548d // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
)

// uncomment when debugging timings
// require internal/timer v1.0.0
replace internal/timer => ./internal/timer
