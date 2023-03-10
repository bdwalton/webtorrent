package controllers

import (
	"fmt"
	"sync"

	"github.com/cenkalti/rain/torrent"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/ini.v1"
)

const (
	suffix = ".meta"
)

// srv is the package-level object we can reference to find the
// Torrent client and global config.
var srv *server

// registerPrometheus is a dumb helper to centralize all prometheus
// registrations for the conrollers package. This could possibly
// trigger panics, so must only ever be called during server startup,
// never later.
func registerPrometheus() {
	// This exports the number of currently loaded torrents. They
	// may or may not be active.
	loadedTorrents := prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "webtorrent_torrents_loaded",
			Help: "The number of currently loaded (not necessarily active) torrents.",
		},
		func() float64 { return float64(len(srv.client.ListTorrents())) },
	)

	prometheus.MustRegister(loadedTorrents)
}

type server struct {
	client *torrent.Session
	cfg    *ini.File
	mtx    sync.Mutex
}

func (s *server) autoStartTorrents() bool {
	switch s.cfg.Section("torrent").Key("autostart").String() {
	case "true":
		return true
	default:
		return false
	}
}

func makeTorrentConfig(cfg *ini.File) torrent.Config {
	tcfg := torrent.DefaultConfig
	tcfg.RPCEnabled = false
	tcfg.DataDir = cfg.Section("torrent").Key("datadir").String()

	return tcfg
}

// Init will create the global srv object and populate it with a
// Torrent client. It also handles initializing pre-saved torrents
// from storage.
func Init(cfg *ini.File) error {
	srv = &server{cfg: cfg}

	tc, err := torrent.NewSession(makeTorrentConfig(cfg))
	if err != nil {
		return fmt.Errorf("Error establishing torrent client: %v\n", err)
	}

	srv.client = tc

	registerPrometheus()

	return nil
}
