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
	@rm -rf cows/
	docker create \
		--name pokebuilder \
		pokesay-go:latest
	@docker cp pokebuilder:/tmp/cows/ .
	@tar czf build/cows.tar.gz cows/
	@rm -rf cows/
	@docker rm -f pokebuilder
	@du -sh build/cows.tar.gz

build/bin: build/docker
	docker create --name pokesay pokesay-go:latest
	docker cp pokesay:/usr/local/src/pokesay .
	docker rm pokesay

build/android:
	rm -f go.mod go.sum
	go mod init github.com/tmck-code/pokesay-go
	go get github.com/mitchellh/go-wordwrap
	go build pokesay.go cows/bindata.go

install: build/bin
	cp -v pokesay $(HOME)/bin/

.PHONY: all clean build/docker build/cows build/bin build/android install
