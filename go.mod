module github.com/tmck-code/pokesay-go

go 1.17

require (
	github.com/mitchellh/go-wordwrap v1.0.1
	github.com/schollz/progressbar/v3 v3.8.3
)

// uncomment when debugging timings
// require internal/timer v1.0.0

require google.golang.org/protobuf v1.27.1 // indirect

replace internal/timer => ./internal/timer
