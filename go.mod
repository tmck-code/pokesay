module github.com/tmck-code/pokesay-go

go 1.17

require github.com/mitchellh/go-wordwrap v1.0.1

require internal/timer v1.0.0

require google.golang.org/protobuf v1.27.1 // indirect

replace internal/timer => ./internal/timer
