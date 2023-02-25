package main

import (
	"flag"
	"fmt"
)

var (
	configFile = flag.String("config", "", "Path to the config file. Required.")
)

func main() {
	flag.Parse()
	fmt.Printf("Config file: %s\n", *configFile)
}
