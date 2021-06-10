VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`
COMMIT=`git rev-list -1 HEAD`
BINARY=axwgameboy

LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} -X main.date=${BUILD} -X main.commit=${COMMIT}"

build:
	echo "Building for linux and windows"
	GOOS=linux GOARCH=amd64 go build -o build/${BINARY}-linux-amd64 ${LDFLAGS} go.janniklasrichter.de/axwgameboy/cmd/axwgameboy
	GOOS=windows GOARCH=amd64 go build -o build/${BINARY}-windows-amd64.exe ${LDFLAGS} go.janniklasrichter.de/axwgameboy/cmd/axwgameboy

run:
	go run go.janniklasrichter.de/axwgameboy/cmd/axwgameboy

test:
	go test -timeout 30s go.janniklasrichter.de/axwgameboy/internal/cpu -v

clean:
	rm -rf build/

all: clean build

.PHONY: clean build run all