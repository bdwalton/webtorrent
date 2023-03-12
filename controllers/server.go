package controllers

import (
	"fmt"
	"path/filepath"
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

func newServer(cfg *ini.File) (*server, error) {
	srv := &server{cfg: cfg}

	tc, err := torrent.NewSession(makeTorrentConfig(cfg))
	if err != nil {
		return nil, fmt.Errorf("Error establishing torrent client: %v\n", err)
	}

	srv.client = tc

	return srv, nil

}

func (s *server) shutdown() {
	s.client.Close()
}

func (s *server) stopOnAdd() bool {
	return s.cfg.Section("torrent").Key("stop_on_add").String() == "true"
}

func (s *server) stopAfterDownload() bool {
	return s.cfg.Section("torrent").Key("stop_on_complete").String() == "true"
}

func (s *server) stopAfterMetadata() bool {
	return s.cfg.Section("torrent").Key("stop_after_metadata").String() == "true"
}

func makeTorrentConfig(cfg *ini.File) torrent.Config {
	tcfg := torrent.DefaultConfig
	tcfg.RPCEnabled = false
	tcfg.DataDir = filepath.Join(cfg.Section("torrent").Key("basedir").String(), "torrents")
	tcfg.Database = filepath.Join(cfg.Section("torrent").Key("basedir").String(), "metadata")

	return tcfg
}

// Init will create the global srv object and populate it with a
// Torrent client. It also handles initializing pre-saved torrents
// from storage.
func Init(cfg *ini.File) error {
	s, err := newServer(cfg)
	if err != nil {
		return err
	}

	srv = s

	registerPrometheus()

	return nil
}

func Shutdown() {
	srv.shutdown()

}
