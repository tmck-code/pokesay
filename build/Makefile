DOCKER_OUTPUT_DIR ?= /tmp/cows
DOCKER_REPO ?= tmckcode
DOCKER_TAG ?= latest

DOCKER_IMAGE=$(DOCKER_REPO)/pokesay:$(DOCKER_TAG)


all: build/docker build/cows build/assets test build/release build/deb

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
		-v $(PWD)/../:/usr/local/src \
		--rm --name pokesay \
	    $(DOCKER_IMAGE) \
			build/scripts/build_assets.sh

build/release: build/assets
	docker run \
		-v $(PWD)/../:/usr/local/src \
		--rm --name pokesay \
		  $(DOCKER_IMAGE) \
		  build/scripts/build.sh
	tree $(PWD)/bin/

build/deb:
	@echo "\n> Building DEB packages"
	@docker run \
		-v $(PWD)/../:/usr/local/src \
		--rm --name pokesay \
		  $(DOCKER_IMAGE) \
		  bash -c "/usr/local/src/build/scripts/build_packages.sh deb"

build/arch:
	@echo "\n> Building Arch packages"
	@docker run \
		-it \
		-v $(PWD)/../:/usr/local/src \
		--platform linux/amd64 \
		--rm --name pokesay \
		  archlinux:base-devel \
		  bash -c "useradd u -m && su - u -c 'bash /usr/local/src/build/scripts/build_packages.sh arch'"
# 		  bash -c "/usr/local/src/build/scripts/build_packages.sh arch_as_user"

build/packages: build/deb build/arch
	@echo "\n> Built packages:"
	@ls -alh $(PWD)/packages/

test:
	docker run \
		-v $(PWD)/../:/usr/local/src \
		--rm --name pokesay-test \
			$(DOCKER_IMAGE) \
			gotestsum --format dots

.PHONY: all build/docker build/cows build/assets build/release test
.PHONY: build/deb build/arch build/packages
