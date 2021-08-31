//go:build !js
// +build !js

package main

import (
	"log"
	"net/http"
)

const (
	addr = "0.0.0.0:8008"
	dir  = ""
)

func main() {
	log.Fatal(http.ListenAndServe(addr, http.FileServer(http.Dir(dir))))
}
