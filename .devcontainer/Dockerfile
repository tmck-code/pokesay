FROM github.com/tmck-code/pokesay:latest

WORKDIR /go
RUN go install github.com/ramya-rao-a/go-outline@latest \
    && go install golang.org/x/tools/gopls@latest \
    && go install github.com/766b/go-outliner@latest
