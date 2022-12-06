
BIN_NAME = zerto-exporter
DOCKER_IMAGE_NAME ?= claranet/zerto-exporter
GOPATH = $($pwd)

all: linux darwin windows

linux: prepare
	$(eval GOOS=linux)
	$(eval GOARCH=amd64)
	go build -o ./bin/$(BIN_NAME)
	zip ./bin/$(BIN_NAME)-$(GOOS)-$(GOARCH).zip ./bin/$(BIN_NAME)

darwin: prepare
	$(eval GOOS=darwin)
	$(eval GOARCH=amd64)
	go build -o ./bin/$(BIN_NAME)
	zip ./bin/$(BIN_NAME)-$(GOOS)-$(GOARCH).zip ./bin/$(BIN_NAME)

windows: prepare
	$(eval GOOS=windows)
	$(eval GOARCH=amd64)
	go build -o ./bin/$(BIN_NAME).exe
	zip ./bin/$(BIN_NAME)-$(GOOS)-$(GOARCH).zip ./bin/$(BIN_NAME).exe

docker:
	@echo ">> Compile using docker container"
	@docker build -t $(DOCKER_IMAGE_NAME):latest .
	@docker tag $(DOCKER_IMAGE_NAME):latest $(DOCKER_IMAGE_NAME):v1.0.1

prepare:
	@echo "Create output directory ./bin/"
	mkdir -p bin/
	@echo "GO get dependencies"
	go get -d

clean:
	@echo "Clean up"
	go clean
	rm -rf bin/

	

.PHONY: all
