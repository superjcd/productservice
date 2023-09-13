package main

import (
	"flag"

	"github.com/HooYa-Bigdata/productservice/cmd/server"
)

var cfg = flag.String("config", "config/config.yaml", "config file location")

func main() {
	flag.Parse()
	server.Run(*cfg)
}
