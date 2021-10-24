VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`
COMMIT=`git rev-list -1 HEAD`
BINARY=axwgameboy

LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} -X main.date=${BUILD} -X main.commit=${COMMIT}"

build: windows linux wasm android

windows:
	GOOS=windows GOARCH=amd64 \
		go build -o build/${BINARY}-windows-amd64.exe ${LDFLAGS} go.janniklasrichter.de/axwgameboy/cmd/axwgameboy

linux:
	GOOS=linux GOARCH=amd64 \
		go build -o build/${BINARY}-linux-amd64 ${LDFLAGS} go.janniklasrichter.de/axwgameboy/cmd/axwgameboy

wasm:
	GOOS=js GOARCH=wasm \
		go build -o build/wasm/${BINARY}-wasm.wasm ${LDFLAGS} go.janniklasrichter.de/axwgameboy/cmd/axwgameboy
	cp assets/gameframe.html build/wasm/
	cp assets/index.html build/wasm/
	cp assets/style.css build/wasm/
	cp assets/wasm_exec.js build/wasm/

wasmserver:
	GOOS=linux GOARCH=amd64 \
		go build -o build/wasm/wasmserver ${LDFLAGS} go.janniklasrichter.de/axwgameboy/cmd/wasmserver
	chmod +x build/wasm/wasmserver

runwasm: wasm wasmserver
	cd build/wasm/
	./wasmserver

android:
	gomobile build -target=android go.janniklasrichter.de/axwgameboy/cmd/axwgameboy

run:
	go run go.janniklasrichter.de/axwgameboy/cmd/axwgameboy $(ARGS)

runwindows: clean build
	./build/${BINARY}-windows-amd64.exe $(ARGS)

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

all: clean build

.PHONY: clean build run all