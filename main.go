package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/anacrolix/torrent"
	"github.com/bdwalton/webtorrent/server"
	"gopkg.in/ini.v1"
)

var (
	configFile  = flag.String("config", "", "Path to the config file. Required.")
	printConfig = flag.Bool("print_config", false, "Print the final config and exit.")
)

//go:embed default.ini
var defaultConfig []byte

// validateConfig expects a valid ini.File object and ensure that all
// of the required settings are valid. A useful error will be returned
// for any object that doesn't meet the basic requirements.
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

	// defaultConfig first because others take precedence when it
	// comes to duplicates.
	cfg, err := ini.Load(defaultConfig, *configFile)
	if err != nil {
		log.Fatalf("Failed to load config %q: %v\n", *configFile, err)
	}

	if err := validateConfig(cfg); err != nil {
		log.Fatalf("Invalid config: %v\n", err)
	}

	if *printConfig {
		if _, err := cfg.WriteTo(os.Stdout); err != nil {
			log.Fatalf("Error printing config: %v\n", err)
		}
	}

	tcfg := torrent.NewDefaultClientConfig()
	c, err := torrent.NewClient(tcfg)
	if err != nil {
		log.Fatalf("Error establishing torrent client: %v\n", err)
	}

	if err := server.ListenAndServe(context.Background(), c, cfg); err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}
