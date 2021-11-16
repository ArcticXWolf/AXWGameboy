VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`
COMMIT=`git rev-list -1 HEAD`
BINARY=axwgameboy

LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.date=${BUILD} -X main.commit=${COMMIT}"

build:
	GOOS=js GOARCH=wasm \
		go build -o build/${BINARY}-wasm.wasm ${LDFLAGS} go.janniklasrichter.de/axwgameboy/cmd/axwgameboy
	cp -r assets/* build/

serve: build
	GOOS=linux GOARCH=amd64 \
		go run go.janniklasrichter.de/axwgameboy/cmd/wasmserver --directory='./build'

test:
	go test go.janniklasrichter.de/axwgameboy/internal

testverbose:
	go test go.janniklasrichter.de/axwgameboy/internal -v

mooneye:
	go test go.janniklasrichter.de/axwgameboy/internal -v -run ^TestMooneyeRoms$

blargg:
	go test go.janniklasrichter.de/axwgameboy/internal -v -run ^TestBlargg.*$

clean:
	rm -rf build/

all: clean build serve

.PHONY: clean build all

# Sameboy bootroms

BIN := build/bin
OBJ := build/obj

ifeq ($(PLATFORM),windows32)
# To force use of the Unix version instead of the Windows version
MKDIR := $(shell which mkdir)
NULL := NUL
_ := $(shell chcp 65001)
EXESUFFIX:=.exe
NATIVE_CC = clang -IWindows -Wno-deprecated-declarations --target=i386-pc-windows
else
MKDIR := mkdir
NULL := /dev/null
EXESUFFIX:=
NATIVE_CC := cc
endif

PB12_COMPRESS := build/pb12$(EXESUFFIX)

bootroms: $(BIN)/bootroms/cgb_boot.bin $(BIN)/bootroms/dmg_boot.bin

$(OBJ)/bootroms/AXWGameboyLogo.pb12: $(OBJ)/bootroms/AXWGameboyLogo.2bpp $(PB12_COMPRESS)
	$(realpath $(PB12_COMPRESS)) < $< > $@

$(OBJ)/bootroms/AXWGameboyLogo.2bpp: internal/bootroms/AXWGameboyLogo.png
	-@$(MKDIR) -p $(dir $@)
	rgbgfx -h -u -o $@ $<
	
$(PB12_COMPRESS): internal/bootroms/pb12.c
	$(NATIVE_CC) -std=c99 -Wall -Werror $< -o $@

$(BIN)/bootroms/%.bin: internal/bootroms/%.asm $(OBJ)/bootroms/AXWGameboyLogo.pb12
	-@$(MKDIR) -p $(dir $@)
	rgbasm -i $(OBJ)/bootroms/ -i internal/bootroms/ -o $@.tmp $<
	rgblink -o $@.tmp2 $@.tmp
	dd if=$@.tmp2 of=$@ count=1 bs=$(if $(findstring dmg,$@)$(findstring sgb,$@),256,2304) 2> $(NULL)
	@rm $@.tmp $@.tmp2
