package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
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
	//go:embed ui/dist/ui
	angularFS embed.FS
)

// validateConfig expects a valid ini.File object and ensure that all
// of the required settings are valid. A useful error will be returned
// for any object that doesn't meet the basic requirements.
func validateConfig(cfg *ini.File) error {
	if _, err := cfg.Section("server").Key("port").Int(); err != nil {
		return fmt.Errorf("Invalid server.port setting: %w", err)
	}

	// torrent.basedir is a required setting so the torrent
	// library knows where to store downloads and metadata.
	bd := cfg.Section("torrent").Key("basedir").String()
	if bd == "" {
		return fmt.Errorf("Invalid torrent.basedir setting: %q", bd)
	}

	if err := os.MkdirAll(bd, 0755); err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("Invalid torrent.basedir setting: %v", err)
	}

	fdd := cfg.Section("torrent").Key("final_datadir").String()
	if fdd == "" || fdd == bd {
		return fmt.Errorf("Must set torrent.final_datadir to something other than torrent.basedir")
	} else {
		if err := os.MkdirAll(fdd, 0755); err != nil && !errors.Is(err, os.ErrExist) {
			return fmt.Errorf("Invalid torrent.final_datadir setting: %v", err)
		}
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

	// This should initialize the torrent client with appropriate config.
	if err := controllers.Init(cfg); err != nil {
		log.Fatalf("WebTorrent: %v", err)
	}

	mappings.Init(cfg, "./ui/dist/ui", fs.FS(angularFS))

	srv := &http.Server{
		Addr:    ":" + cfg.Section("server").Key("port").String(),
		Handler: mappings.GetRouter(),
	}

	go func() {
		log.Println("WebTorrent: Serving on port", cfg.Section("server").Key("port").String())
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("WebTorrent: Server error: %v\n", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-sig:
		log.Printf("WebTorrent: Received signal %q. Initiating shutdown.", s)
		controllers.Shutdown()
	}

	sctx, release := context.WithTimeout(context.Background(), 10*time.Second)
	defer release()
	srv.Shutdown(sctx)

	log.Println("WebTorrent: Goodbye...")
}
