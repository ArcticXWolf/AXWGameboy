//go:build !js
// +build !js

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	addr string
	port int
	dir  string
)

func init() {
	flag.StringVar(&dir, "directory", ".", "directory to serve")
	flag.StringVar(&addr, "bindip", "0.0.0.0", "ip to bind onto")
	flag.IntVar(&port, "bindport", 8008, "port to bind onto")
}

func main() {
	flag.Parse()
	log.Printf("Serving %s on %s:%d", dir, addr, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), http.FileServer(http.Dir(dir))))
}
