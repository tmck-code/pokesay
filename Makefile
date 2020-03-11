build-docker:
	docker build -f ops/Dockerfile -t pokesay-go:latest .

build:
	docker run -it --name pokebuilder pokesay-go:latest bash -c "/usr/local/src/build.sh /usr/local/src/icons /tmp"
	@docker cp pokebuilder:/tmp/cows/ .
	@tar czf cows.tar.gz cows/
	@rm -rf cows/
	@docker rm -f pokebuilder

install:
	@./install.sh

.PHONY: build-docker build install
