package main

import (
	"flag"
	"fmt"
	"log"

	"gopkg.in/ini.v1"
)

var (
	configFile = flag.String("config", "", "Path to the config file. Required.")
)

func validateConfig(cfg *ini.File) error {
	if _, err := cfg.Section("server").Key("port").Int(); err != nil {
		return fmt.Errorf("Invalid server.port setting: %w", err)
	}

	return nil
}

func main() {
	flag.Parse()

	if *configFile == "" {
		log.Fatal("You must specify a config file via --config.")
	}

	cfg, err := ini.Load(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config %q: %v\n", *configFile, err)
	}

	if err := validateConfig(cfg); err != nil {
		log.Fatalf("Invalid config: %v\n", err)
	}
}
