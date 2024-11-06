DOCKER_OUTPUT_DIR ?= /tmp/cows
DOCKER_REPO ?= tmckcode
DOCKER_TAG ?= latest

DOCKER_IMAGE=$(DOCKER_REPO)/pokesay:$(DOCKER_TAG)


all: build/docker build/cows build/assets test build/release

build/docker:
	echo "Building $(DOCKER_IMAGE)"
	docker build -f Dockerfile -t $(DOCKER_IMAGE) ..

build/cows:
	@rm -rf cows.tar.gz cows/
	docker rm -f pokebuilder
	docker create --name pokebuilder $(DOCKER_IMAGE)
	@docker cp pokebuilder:$(DOCKER_OUTPUT_DIR)/ cows/
	@tar czf cows.tar.gz cows/
	@rm -rf cows/
	@docker rm -f pokebuilder
	@du -sh cows.tar.gz

# generate embedded bin files for category/metadata/the actual pokemon
build/assets:
	docker run \
		-v $(PWD)/assets:/usr/local/src/build/assets \
		--rm --name pokesay \
	    $(DOCKER_IMAGE) \
			build/scripts/build_assets.sh

build/release: build/assets
	docker run \
		-v $(PWD)/bin:/usr/local/src/build/bin \
		-v $(PWD)/assets:/usr/local/src/build/assets \
		--rm --name pokesay \
		  $(DOCKER_IMAGE) \
		  build/scripts/build.sh
	tree -s $(PWD)/bin/

test:
	docker run \
		-v $(PWD)/../:/usr/local/src \
		--rm --name pokesay-test \
			$(DOCKER_IMAGE) \
			go test -v ./test/

.PHONY: all build/docker build/cows build/assets build/release test
