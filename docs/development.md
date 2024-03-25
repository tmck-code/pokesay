# Developing

## Testing


```shell
go test -v ./test/

# or, with debug information printed
DEBUG=test go test -v ./test/
```

## Deploying

1. Make a PR, and merge it into master
2. Draft a new release [https://github.com/tmck-code/pokesay/releases/new](https://github.com/tmck-code/pokesay/releases/new)
3. Give it a new tag and identical release title by incrementing the version
   1. Generate the release notes automatically
4. Run `make all` (or some variation) locally to create the binaries under `build/bin/`
   1. Upload the bins to the release page
5. Then click "Publish Release"! The install script downloads binaries from the most recently-published released
