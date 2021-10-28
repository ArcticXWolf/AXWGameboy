VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`
COMMIT=`git rev-list -1 HEAD`
BINARY=axwgameboy

LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} -X main.date=${BUILD} -X main.commit=${COMMIT}"

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