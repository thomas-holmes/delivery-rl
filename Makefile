GO_BIN?=go
CC=$(GO_BIN)
GOBUILD=$(CC) build
GOPKG='github.com/thomas-holmes/delivery-rl/game'
MKFILE_PATH=$(abspath $(lastword $(MAKEFILE_LIST)))
ROOT=$(shell git rev-parse --show-toplevel)
SHA=$(shell git rev-parse HEAD)
VERSION=$(shell git describe --tags --always)
BUILDPATH=dist/deliveryrl-$(VERSION)
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

binclean:
	rm run_dir/$(BINARY) | true
	rm run_dir/$(BINARY).exe | true

distclean:
	rm -rf dist

dist: binclean distclean test build buildwin
	mkdir -p $(BUILDPATH)
	echo $(SHA) > $(BUILDPATH)/COMMIT
	cp README.md $(BUILDPATH)
	cp CHANGELOG.md $(BUILDPATH)
	cp -r run_dir/* $(BUILDPATH)
	cd dist && zip -r ../deliveryrl-$(VERSION).zip deliveryrl-$(VERSION)

test:
	$(GO_BIN) test ./...
