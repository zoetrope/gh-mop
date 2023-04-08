
.PHONY: all
all: build

build:
	@echo "Building..."
	@go build main.go

test:
	@echo "Testing..."
	@go test -v -count 1 ./...
