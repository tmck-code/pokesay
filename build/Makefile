DOCKER_OUTPUT_DIR ?= /tmp/cows
DOCKER_REPO ?= tmckcode
DOCKER_TAG ?= latest

DOCKER_IMAGE=$(DOCKER_REPO)/pokesay:$(DOCKER_TAG)
DEBUG ?= ""

ifeq ($(shell which gecho > /dev/null 2>&1 && echo 1 || echo 0), 1)
	echo := gecho
else
	echo := echo
endif

all: build/docker build/cows build/assets test build/bin build/packages

clean:
	@$(echo) -e "\e[48;5;30m> Cleaning up\e[0m"
	rm -rf ../dist/ assets/ deb/ arch/ cows/ bin/

build/docker:
	@$(echo) -e "\e[48;5;30m> Building docker image $(DOCKER_IMAGE)\e[0m"
	@docker build \
		--platform linux/amd64 \
		-f Dockerfile \
		-t $(DOCKER_IMAGE) \
		--build-arg DEBUG=$(DEBUG) \
		..

build/docker-convert:
	@$(echo) -e "\e[48;5;30m> Building docker image $(DOCKER_IMAGE)-convert\e[0m"
	@docker build \
		--platform linux/amd64 \
		-f Dockerfile.convert \
		-t $(DOCKER_REPO)/pokesay-convert:$(DOCKER_TAG) \
		--build-arg DEBUG=$(DEBUG) \
		..

build/cows:
	@$(echo) -e "\e[48;5;30m> Building cows\e[0m"
	@rm -rf cows.tar.gz cows/
	@mkdir -p cows/
	@docker run --rm \
		--platform linux/amd64 \
		-v $(PWD)/../:/usr/local/src \
		-v $(PWD)/pokesprite:/tmp/original/pokesprite \
		-v $(PWD)/cows:/tmp/cows \
		-e DEBUG=$(DEBUG) \
		-u u \
		$(DOCKER_IMAGE) \
		  build/scripts/build_cowfiles.sh
	@tar -czf cows.tar.gz ./cows
	@du -sh cows.tar.gz

# generate embedded bin files for category/metadata/the actual pokemon
build/assets:
	@$(echo) -e "\e[48;5;30m> Building assets\e[0m"
	@mkdir -p $(PWD)/assets
	@rm -rf $(PWD)/cows/
	@tar -xzf cows.tar.gz
	@docker run --rm \
		-v $(PWD)/../:/usr/local/src \
		-v $(PWD)/cows:/tmp/cows \
		--platform linux/amd64 \
		-u u \
		--name pokesay \
	    $(DOCKER_IMAGE) \
			build/scripts/build_assets.sh
	@tree -L 1 $(PWD)/assets/

build/bin:
	@$(echo) -e "\e[48;5;30m> Building binaries\e[0m"
	@mkdir -p $(PWD)/../dist/bin/ $(PWD)/../dist/packages/ $(PWD)/../dist/tarballs/
	@docker run --rm \
		-v $(PWD)/../:/usr/local/src \
		-e VERSION=$(VERSION) \
		-e DEBUG=$(DEBUG) \
		--platform linux/amd64 \
		--user u \
		--name pokesay \
		  $(DOCKER_IMAGE) \
		  /usr/local/src/build/scripts/build_bin.sh
	@tree $(PWD)/../dist/

build/deb:
	@$(echo) -e "\e[48;5;30m> Building DEB packages\e[0m"
	@docker run --rm \
		-v $(PWD)/../:/usr/local/src \
		-e VERSION=$(VERSION) \
		--platform linux/amd64 \
		--name pokesay \
		  $(DOCKER_IMAGE) \
		  bash -c "/usr/local/src/build/scripts/build_packages.sh deb"

build/arch:
	@$(echo) -e "\e[48;5;30m> Building ARCH packages\e[0m"
	@docker run --rm\
		-v $(PWD)/../:/usr/local/src \
		-e VERSION=$(VERSION) \
		--platform linux/amd64 \
		--name pokesay \
		  archlinux:base-devel \
		  bash -c "useradd u -m && VERSION=$(VERSION) /usr/local/src/build/scripts/build_packages.sh arch"

build/packages: build/deb build/arch
	@$(echo) -e "\e[48;5;30m> Built packages:\e[0m"
	@tree $(PWD)/../dist/

test:
	@$(echo) -e "\e[48;5;30m> Running tests\e[0m"
	@docker run --rm \
		-v $(PWD)/../:/usr/local/src \
		--platform linux/amd64 \
		--user root \
		--name pokesay-test \
			$(DOCKER_IMAGE) \
			gotestsum --format dots

.PHONY: all clean build/docker build/cows build/assets build/bin test
.PHONY: build/deb build/arch build/packages
