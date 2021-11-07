# pokesay-go
A (much) faster go version of tmck-code/pokesay

## One-line install

```shell
# OSX / darwin
bash -c "$(curl https://raw.githubusercontent.com/tmck-code/pokesay-go/master/install.sh)" bash darwin amd64

# Linux x64
bash -c "$(curl https://raw.githubusercontent.com/tmck-code/pokesay-go/master/install.sh)" bash linux amd64

# Android ARM64 (Termux)
bash -c "$(curl https://raw.githubusercontent.com/tmck-code/pokesay-go/master/install.sh)" bash android arm64
```

## Building binaries

```shell
TARGET_GOOS=darwin TARGET_GOARCH=amd64 make build/bin
TARGET_GOOS=linux TARGET_GOARCH=amd64 make build/bin
TARGET_GOOS=windows TARGET_GOARCH=amd64 make build/bin
TARGET_GOOS=android TARGET_GOARCH=arm64 make build/bin
```
