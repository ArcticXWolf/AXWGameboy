package viewcontainer

import (
	"encoding/json"
	"io"
	"syscall/js"
)

type localStorageSavegame struct {
	romHash    string
	cacheData  []byte
	readOffset int64
}

func (lSS *localStorageSavegame) Read(p []byte) (n int, err error) {
	if lSS.cacheData == nil {
		item := js.Global().Get("localStorage").Call("getItem", lSS.romHash)
		if !item.Truthy() {
			err = io.EOF
			return
		}
		data := []byte(item.String())
		err = json.Unmarshal(data, &lSS.cacheData)
		if err != nil {
			return
		}
	}

	if lSS.readOffset >= int64(len(lSS.cacheData)) {
		err = io.EOF
		return
	}

	n = copy(p, lSS.cacheData[lSS.readOffset:])
	lSS.readOffset += int64(n)
	return
}

func (lSS *localStorageSavegame) Write(data []byte) (n int, err error) {
	storage := js.Global().Get("localStorage")

	marshalledData, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	storage.Call("setItem", lSS.romHash, string(marshalledData))
	return len(data), err
}
