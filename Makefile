TARGET_GOOS       ?= linux
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
		--build-arg GOOS=$(TARGET_GOOS) \
		--build-arg GOARCH=$(TARGET_GOARCH) \
		-f build/Dockerfile \
		-t pokesay-go:latest .

build/cows:
	@rm -rf build/cows/ cmd/bindata.go
	docker create \
		--name pokebuilder \
		pokesay-go:latest
	@docker cp pokebuilder:/tmp/cows/ build/
	@docker cp pokebuilder:$(DOCKER_BUILD_DIR)/cmd/bindata.go cmd/
	@tar czf build/cows.tar.gz build/cows/
	@rm -rf build/cows/
	@docker rm -f pokebuilder
	@du -sh build/cows.tar.gz

build/bin: build/docker
	docker create --name pokesay pokesay-go:latest
	docker cp pokesay:/usr/local/src/pokesay .
	docker rm pokesay
	mv -v pokesay pokesay-$(TARGET_GOOS)-$(TARGET_GOARCH)

build/android:
	go mod tidy
	go get github.com/mitchellh/go-wordwrap
	go get github.com/go-bindata/go-bindata
	tar xzf build/cows.tar.gz -C build/
	go-bindata -o cmd/bindata.go build/cows/...
	go build cmd/pokesay.go cmd/bindata.go
	mv -v pokesay build/pokesay-android-arm64
	rm -rf build/cows

install:
	cp -v pokesay-$(TARGET_GOOS)-$(TARGET_GOARCH) $(HOME)/bin/pokesay

.PHONY: all clean build/docker build/cows build/bin build/android install
