TARGET_GOOS       ?= linux
TARGET_GOARCH     ?= amd64
DOCKER_BUILD_DIR  ?= /usr/local/src
DOCKER_OUTPUT_DIR ?= /tmp/cows

all: build/docker build/cows build/release

build/docker:
	docker build \
		-f Dockerfile \
		-t pokesay-go:latest ..

build/cows:
	@rm -rf cows.tar.gz cows/
	docker create \
		--name pokebuilder \
		pokesay-go:latest
	@docker cp pokebuilder:$(DOCKER_OUTPUT_DIR)/ cows/
	@docker cp pokebuilder:$(DOCKER_BUILD_DIR)/cmd/bindata.go cmd/
	@tar czf cows.tar.gz cows/
	@rm -rf cows/
	@docker rm -f pokebuilder
	@du -sh cows.tar.gz

install: build/bin
	cp -v pokesay-$(TARGET_GOOS)-$(TARGET_GOARCH) $(HOME)/bin/pokesay

build/release: build/docker
	@docker create --name pokesay pokesay-go:latest
	@docker cp pokesay:/usr/local/src/pokesay-linux-amd64 .
	@docker cp pokesay:/usr/local/src/pokesay-darwin-amd64 .
	@docker cp pokesay:/usr/local/src/pokesay-windows-amd64.exe .
	@docker cp pokesay:/usr/local/src/pokesay-android-arm64 .
	@docker rm -f pokesay

.PHONY: all build/docker build/cows install build/release
