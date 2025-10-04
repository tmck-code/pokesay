DOCKER_OUTPUT_DIR ?= /tmp/cows
DOCKER_REPO ?= tmckcode
DOCKER_TAG ?= latest

DOCKER_IMAGE=$(DOCKER_REPO)/pokesay:$(DOCKER_TAG)
DEBUG ?= ""


all: build/docker build/cows build/assets test build/bin build/packages

build/docker:
	echo "Building $(DOCKER_IMAGE)"
# 	docker build --platform linux/amd64 -f Dockerfile -t $(DOCKER_IMAGE) ..
	docker build -f Dockerfile -t $(DOCKER_IMAGE) ..

build/cows:
	@echo -e "\n> Building cows"
	@rm -rf cows.tar.gz cows/
	@docker inspect pokebuilder > /dev/null 2>&1 && docker rm -f pokebuilder || true
	docker create --name pokebuilder $(DOCKER_IMAGE)
	docker cp pokebuilder:$(DOCKER_OUTPUT_DIR)/ cows/
	@tar czf cows.tar.gz cows/
	@rm -rf cows/
	@docker rm -f pokebuilder
	@du -sh cows.tar.gz

# generate embedded bin files for category/metadata/the actual pokemon
build/assets:
	@docker run \
		-v $(PWD)/../:/home/u/pokesay \
		--rm --name pokesay \
	    $(DOCKER_IMAGE) \
			build/scripts/build_assets.sh
	@tree -L 1 assets

build/bin:
	@echo -e "\n> Building binaries"
	@mkdir -p $(PWD)/../dist/bin/ $(PWD)/../dist/packages/ $(PWD)/../dist/tarballs/
	@docker run \
		-v $(PWD)/../:/home/u/pokesay \
		-e VERSION=$(VERSION) \
		-e DEBUG=$(DEBUG) \
		--rm --name pokesay \
		  $(DOCKER_IMAGE) \
		  build/scripts/build_bin.sh
	@tree $(PWD)/../dist/

build/deb:
	@echo -e "\n> Building DEB packages"
	@docker run --rm \
		-v $(PWD)/../:/home/u/pokesay \
		-e VERSION=$(VERSION) \
		--name pokesay \
		  $(DOCKER_IMAGE) \
		  bash -c "build/scripts/build_packages.sh deb"

build/arch:
	@echo -e "\n> Building Arch packages"
	@docker run --rm \
		-v $(PWD)/../:/home/u/pokesay \
		-e VERSION=$(VERSION) \
		--platform linux/amd64 \
		--name pokesay \
		  archlinux:base-devel \
		  bash -c "useradd u && /home/u/pokesay/build/scripts/build_packages.sh arch"

build/packages: build/deb build/arch
	@echo -e "\n> Built packages:"
	@tree $(PWD)/../dist/

test:
	@echo -e "\n> Running tests"
	@docker run \
		-v $(PWD)/../:/home/u/pokesay \
		--rm --name pokesay-test \
			$(DOCKER_IMAGE) \
			gotestsum --format dots

.PHONY: all build/docker build/cows build/assets build/bin test
.PHONY: build/deb build/arch build/packages
