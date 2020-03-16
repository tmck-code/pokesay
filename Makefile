TARGET_GOOS 				?= "darwin"
TARGET_GOARCH 			?= "amd64"
DOCKER_BUILD_DIR  ?= "/usr/local/src"
DOCKER_OUTPUT_DIR ?= "/tmp"


all: clean build/docker build/cows build/bin install

clean:
	rm -rf cows/
	rm -f cowsay
	docker rm -f pokebuilder || echo 'no container to remove'

build/docker:
	docker build -f ops/Dockerfile -t pokesay-go:latest .

build/cows:
	docker run -it \
		--name pokebuilder \
		-e DOCKER_BUILD_DIR=$(DOCKER_BUILD_DIR) \
		-e DOCKER_OUTPUT_DIR=$(DOCKER_OUTPUT_DIR) \
		pokesay-go:latest \
		bash -c "/usr/local/src/build_cows.sh"
	@docker cp pokebuilder:/tmp/cows/ .
	@tar czf cows.tar.gz cows/
	@rm -rf cows/
	@docker rm -f pokebuilder

build/bin:
	docker run -it \
		-v $(PWD)/cows/:/cows \
		--name pokebuilder \
		-e DOCKER_BUILD_DIR=$(DOCKER_BUILD_DIR)/go-cowsay \
		-e DOCKER_OUTPUT_DIR=$(DOCKER_OUTPUT_DIR) \
		-e TARGET_GOOS=$(TARGET_GOOS) \
		-e TARGET_GOARCH=$(TARGET_GOARCH) \
		pokesay-go:latest \
		bash -c "/usr/local/src/build_bin.sh"
	@docker cp pokebuilder:/tmp/cowsay .
	@docker rm -f pokebuilder

build/android:
	go get -u -v github.com/msmith491/go-cowsay || true
	cd $(GOPATH)/src/github.com/msmith491/go-cowsay; \
		make
	cp -v $(GOPATH)/src/github.com/msmith491/go-cowsay/cowsay .

install:
	@./install.sh

.PHONY: all clean build/docker build/cows build/bin install
