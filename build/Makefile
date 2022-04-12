TARGET_GOOS       ?= linux
TARGET_GOARCH     ?= amd64
DOCKER_BUILD_DIR  ?= /usr/local/src
DOCKER_OUTPUT_DIR ?= /tmp/cows

all: build/docker build/cows build/release

build/docker:
	docker build --no-cache -f Dockerfile -t pokesay-go:latest ..

build/cows:
	@rm -rf cows.tar.gz cows/
	docker rm -f --name pokebuilder
	docker create --name pokebuilder pokesay-go:latest
	@docker cp pokebuilder:$(DOCKER_OUTPUT_DIR)/ cows/
	@tar czf cows.tar.gz cows/
	@rm -rf cows/
	@docker rm -f pokebuilder
	@du -sh cows.tar.gz

build/release:
	docker run \
		-v $(PWD)/bin:/usr/local/src/build/bin \
		--rm --name pokesay pokesay-go:latest build/build.sh
	tree

.PHONY: all build/docker build/cows install build/release
