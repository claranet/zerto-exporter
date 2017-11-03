
BIN_NAME = zerto-exporter
DOCKER_IMAGE_NAME ?= zerto-exporter
GOPATH = $($pwd)

all: build

build:
	@echo "Create output directory ./bin/"
	mkdir -p bin/
	@echo "GO get dependencies"
	go get -d
	@echo "Build ..."
	go build -o ./bin/$(BIN_NAME)

clean:
	@echo "Clean up"
	go clean
	rm -rf bin/

docker:
	@echo ">> Compile using docker container"
	@docker build -t "$(DOCKER_IMAGE_NAME)" .

.PHONY: all
