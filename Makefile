CC=go
GOBUILD=$(CC) build
GOPKG='github.com/thomas-holmes/delivery-rl/game'
MKFILE_PATH=$(abspath $(lastword $(MAKEFILE_LIST)))
ROOT=$(notdir $(patsubst %/,%,$(dir $(MKFILE_PATH))))
BINARY=delivery-rl

all: test build

build:
	$(GOBUILD) -o run_dir/$(BINARY) github.com/thomas-holmes/delivery-rl/game

run: build
	cd ./run_dir && ./delivery-rl -no-vsync
