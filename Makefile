TARGET_GOOS       ?= darwin
TARGET_GOARCH     ?= amd64
DOCKER_BUILD_DIR  ?= /usr/local/src
DOCKER_OUTPUT_DIR ?= /tmp
COWFILE_BUILD_DIR ?= /cows

all: clean build/docker build/cows build/bin

clean:
	rm -rf cows/
	docker rm -f pokebuilder || echo 'no container to remove'

build/docker:
	docker build \
		-f Dockerfile \
		-t pokesay-go:latest .

build/cows:
	docker run -it \
		--name pokebuilder \
		-e DOCKER_BUILD_DIR=$(DOCKER_BUILD_DIR) \
		-e DOCKER_OUTPUT_DIR=$(DOCKER_OUTPUT_DIR) \
		pokesay-go:latest \
		bash -c "/usr/local/src/build_cows.sh"
	@rm -rf cows/
	@docker cp pokebuilder:/tmp/cows/ .
	@tar czf cows.tar.gz cows/
	@rm -rf cows/
	@docker rm -f pokebuilder
	@du -sh cows.tar.gz

build/bin: build/docker
	docker create --name pokesay pokesay-go:latest
	docker cp pokesay:/usr/local/src/pokesay .
	docker rm pokesay

build/android:
	rm -f go.mod go.sum
	go mod init github.com/tmck-code/pokesay-go
	go install github.com/go-bindata/go-bindata/...@latest
	go get github.com/mitchellh/go-wordwrap

install: build/bin
	cp -v pokesay $(HOME)/bin/

.PHONY: all clean build/docker build/cows build/bin build/android install
