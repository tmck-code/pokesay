# pokesay-go
A (much) faster go version of tmck-code/pokesay

## Building binaries

```shell
TARGET_GOOS=darwin TARGET_GOARCH=amd64 make build/bin
TARGET_GOOS=linux TARGET_GOARCH=amd64 make build/bin
TARGET_GOOS=windows TARGET_GOARCH=amd64 make build/bin
TARGET_GOOS=android TARGET_GOARCH=arm64 make build/bin
```
