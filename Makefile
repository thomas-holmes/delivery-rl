GO_BIN?=go
CC=$(GO_BIN)
GOBUILD=$(CC) build
GOPKG='github.com/thomas-holmes/delivery-rl/game'
MKFILE_PATH=$(abspath $(lastword $(MAKEFILE_LIST)))
ROOT=$(shell git rev-parse --show-toplevel)
SHA=$(shell git rev-parse HEAD)
VERSION=$(shell git describe --tags --always)
BINARY=delivery-rl

all: test build

build:
	$(GOBUILD) -o run_dir/$(BINARY) github.com/thomas-holmes/delivery-rl/game

buildwin:
	CGO_ENABLED="1" CC="/usr/bin/x86_64-w64-mingw32-gcc" GOOS="windows" CGO_LDFLAGS="-lmingw32 -lSDL2 -I $(ROOT)/sdl2/win/SDL2-2.0.8/x86_64-w64-mingw32/include -L $(ROOT)/sdl2/win/SDL2-2.0.8/x86_64-w64-mingw32/lib" CGO_CFLAGS="-D_REENTRANT -I $(ROOT)/sdl2/win/SDL2-2.0.8/x86_64-w64-mingw32/include -L $(ROOT)/sdl2/win/SDL2-2.0.8/x86_64-w64-mingw32/lib" $(GOBUILD) -o run_dir/$(BINARY).exe github.com/thomas-holmes/delivery-rl/game

run: build
	cd ./run_dir && ./delivery-rl

runv: build
	cd ./run_dir && ./delivery-rl -no-vsync

runrv: build
	cd ./run_dir && ./delivery-rl -no-vsync -seed 0xDEADBEEF

distclean:
	rm run_dir/$(BINARY) || true
	rm run_dir_$(BINARY).exe || true

dist: distclean test build buildwin
	echo $(SHA) > run_dir/COMMIT
	cp README.md run_dir/
	mv run_dir deliveryrl
	tar -czf deliveryrl-$(VERSION)-7drl.tgz deliveryrl
	mv deliveryrl run_dir
	rm run_dir/COMMIT
	rm run_dir/README.md

test:
	$(GO_BIN) test ./...
