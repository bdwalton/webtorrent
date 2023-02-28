package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bdwalton/webtorrent/controllers"
	"github.com/bdwalton/webtorrent/mappings"
	"gopkg.in/ini.v1"
)

var (
	configFile = flag.String("config", "", "Path to the config file. Required.")
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

	// torrent.datadir is a required setting so the torrent
	// library knows where to store download.
	dd := cfg.Section("torrent").Key("datadir").String()
	if dd == "" {
		return fmt.Errorf("Invalid torrent.datadir setting: %q", dd)
	}

	if err := os.MkdirAll(dd, 0755); err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("Invalid torrent.datadir setting: %v", err)
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

	// This should initialize the torrent client with appropriate config.
	if err := controllers.Init(cfg); err != nil {
		log.Fatalf("TorrenServer: %v")
	}

	mappings.Init()

	srv := &http.Server{
		Addr:    ":" + cfg.Section("server").Key("port").String(),
		Handler: mappings.GetRouter(),
	}

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("TorrentServer: Server error: %v\n", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-sig:
		log.Printf("TorrentServer: Received signal %q. Initiating shutdown.", s)
		controllers.ShutdownTorrentClient()
	}

	sctx, release := context.WithTimeout(context.Background(), 10*time.Second)
	defer release()
	srv.Shutdown(sctx)

	log.Println("TorrentServer: Goodbye...")
}
