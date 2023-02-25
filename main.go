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

	ctx := context.Background()

	srv := &http.Server{
		Addr: ":" + cfg.Section("server").Key("port").String(),
	}

	go func() {
		if err := server.ListenAndServe(ctx, srv, c, cfg); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("TorrentServer: Server error: %v\n", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sig:
		log.Println("TorrentServer: Shutting down after signal.")
		c.Close()
	case <-c.Closed():
		log.Println("TorrentServer: Shutting down after torrent client stopped.")

	}

	sctx, release := context.WithTimeout(ctx, 10*time.Second)
	defer release()
	srv.Shutdown(sctx)

	log.Println("TorrentServer: Goodbye...")
}
