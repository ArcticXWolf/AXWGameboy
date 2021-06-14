VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`
COMMIT=`git rev-list -1 HEAD`
BINARY=axwgameboy

LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} -X main.date=${BUILD} -X main.commit=${COMMIT}"

build:
	echo "Building for linux and windows"
	GOOS=linux GOARCH=amd64 \
		go build -o build/${BINARY}-linux-amd64 ${LDFLAGS} go.janniklasrichter.de/axwgameboy/cmd/axwgameboy
	GOOS=windows GOARCH=amd64 \
		CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc \
		go build -o build/${BINARY}-windows-amd64.exe ${LDFLAGS} go.janniklasrichter.de/axwgameboy/cmd/axwgameboy

run:
	go run go.janniklasrichter.de/axwgameboy/cmd/axwgameboy

runwindows: clean build
	./build/${BINARY}-windows-amd64.exe

test:
	go test go.janniklasrichter.de/axwgameboy/internal -v

clean:
	rm -rf build/

all: clean build

.PHONY: clean build run all