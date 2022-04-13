TARGET_GOOS       ?= linux
TARGET_GOARCH     ?= amd64
DOCKER_BUILD_DIR  ?= /usr/local/src
DOCKER_OUTPUT_DIR ?= /tmp/cows

all: build/docker build/cows build/release

build/docker:
	docker build -f Dockerfile -t pokesay-go:latest ..

build/cows:
	@rm -rf cows.tar.gz cows/
	docker rm -f --name pokebuilder
	docker create --name pokebuilder pokesay-go:latest
	@docker cp pokebuilder:$(DOCKER_OUTPUT_DIR)/ cows/
	@tar czf cows.tar.gz cows/
	@rm -rf cows/
	@docker rm -f pokebuilder
	@du -sh cows.tar.gz

# generate embedded bin files for category/metadata/the actual pokemon
build/assets:
	docker run \
		-v $(PWD)/assets:/usr/local/src/build/assets \
		--rm --name pokesay pokesay-go:latest \
			build/build_assets.sh

build/release:
	docker run \
		-v $(PWD)/bin:/usr/local/src/build/bin \
		-v $(PWD)/assets:/usr/local/src/build/assets \
		--rm --name pokesay pokesay-go:latest build/build.sh
	tree bin/

.PHONY: all build/docker build/cows install build/release
