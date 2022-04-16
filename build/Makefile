DOCKER_OUTPUT_DIR ?= /tmp/cows
DOCKER_IMAGE ?= "ghcr.io/tmck-code"
DOCKER_TAG ?= "latest"

all: build/docker build/cows build/assets build/release

build/docker:
	docker build -f Dockerfile -t $(DOCKER_IMAGE)/pokesay-go:$(DOCKER_TAG) ..

build/cows:
	@rm -rf cows.tar.gz cows/
	docker rm -f pokebuilder
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

build/release: build/assets
	docker run \
		-v $(PWD)/bin:/usr/local/src/build/bin \
		-v $(PWD)/assets:/usr/local/src/build/assets \
		--rm --name pokesay pokesay-go:latest build/build.sh
	tree bin/

.PHONY: all build/docker build/cows build/assets build/release
