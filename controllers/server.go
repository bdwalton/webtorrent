package controllers

import (
	"fmt"
	"log"
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
	doneC  chan struct{}
	wg     sync.WaitGroup // used to ensure all watchTorrent goroutines shut down
}

func newServer(cfg *ini.File) (*server, error) {
	srv := &server{
		cfg:   cfg,
		doneC: make(chan struct{}),
	}

	tc, err := torrent.NewSession(makeTorrentConfig(cfg))
	if err != nil {
		return nil, fmt.Errorf("Error establishing torrent client: %v\n", err)
	}

	srv.client = tc

	return srv, nil

}

// watchTorrents should be called at server startup. It will fire up
// monitors to watch existing, running torrents for important state
// changes.
func (s *server) watchTorrents() {
	for _, t := range s.client.ListTorrents() {
		if t.Stats().Status.String() != "Stopped" {
			go s.watchTorrent(t)
		}
	}
}

func (s *server) watchTorrent(t *torrent.Torrent) {
	s.wg.Add(1)

	log.Printf("WebTorrent: watchTorrent(%s) running...", t.ID())

	select {
	case <-s.doneC:
		log.Printf("WebTorrent: watchTorrent(%s) shutting down.", t.ID())
	case <-t.NotifyComplete():
		log.Printf("WebTorrent: watchTorrent(%s) is done.", t.ID())
	case <-t.NotifyStop():
		log.Printf("WebTorrent: watchTorrent(%s) is stopped.", t.ID())
	case <-t.NotifyClose():
		log.Printf("WebTorrent: watchTorrent(%s) was closed (dropped, data removed).", t.ID())
	}

	s.wg.Done()
}

func (s *server) shutdown() {
	s.client.Close()
	close(s.doneC)
	s.wg.Wait()
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

	srv.watchTorrents()

	return nil
}

func Shutdown() {
	srv.shutdown()

}
