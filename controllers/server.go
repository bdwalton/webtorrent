package controllers

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
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

	tc, err := torrent.NewSession(makeTorrentConfig(srv))
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
			s.wg.Add(1)
			go s.watchTorrent(t)
		}
	}
}

func (s *server) watchTorrent(t *torrent.Torrent) {

	log.Printf("WebTorrent: watchTorrent(%s) running...", t.ID())

	select {
	case <-s.doneC:
		log.Printf("WebTorrent: watchTorrent(%s) shutting down.", t.ID())
	case <-t.NotifyComplete():
		log.Printf("WebTorrent: watchTorrent(%s) is done.", t.ID())
		s.wg.Add(1)
		go s.persistTorrent(t)
	case <-t.NotifyStop():
		log.Printf("WebTorrent: watchTorrent(%s) is stopped.", t.ID())
	case <-t.NotifyClose():
		log.Printf("WebTorrent: watchTorrent(%s) was closed (dropped, data removed).", t.ID())
	}

	s.wg.Done()
}

func (s *server) persistTorrent(t *torrent.Torrent) {
	defer s.wg.Done()

	if files, err := t.FilePaths(); err == nil {
		tpd := s.torrentBaseDir()
		fdd := s.finalDataDir()
		if s.datadirIncludesTorrentID() {
			tpd = filepath.Join(tpd, t.ID())
		}

		for _, f := range files {
			src := filepath.Join(tpd, f)
			dst := filepath.Join(fdd, f)
			if err := os.MkdirAll(filepath.Dir(dst), s.filePermissions()); err != nil && !errors.Is(err, os.ErrExist) {
				log.Printf("WebTorrent: persistTorrent(%s) mkdirall error: %v", t.ID(), err)
				return
			}

			if err := os.Link(src, dst); err != nil {
				log.Printf("WebTorrent: persistTorrent(%s) link error: %v", t.ID(), err)
				return
			}
		}

		log.Printf("WebTorrent: persistTorrent(%s) successfully persisted data.", t.ID())
	} else {
		log.Printf("WebTorrent: persistTorrent(%s) error listing files: %v", t.ID(), err)
	}
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

func (s *server) torrentBaseDir() string {
	return filepath.Join(s.cfg.Section("torrent").Key("basedir").String(), "torrents")
}

func (s *server) finalDataDir() string {
	return s.cfg.Section("torrent").Key("final_datadir").String()
}

func (s *server) torrentMetadatadir() string {
	return filepath.Join(s.cfg.Section("torrent").Key("basedir").String(), "metadata")
}

// datadirIncludesTorrentID returns a boolean value representing the
// config key and defaults to true
func (s *server) datadirIncludesTorrentID() bool {
	if s.cfg.Section("torrent").HasKey("datadir_includes_torrentid") {
		return s.cfg.Section("torrent").Key("datadir_includes_torrentid").String() == "true"
	}

	return true
}

func (s *server) filePermissions() fs.FileMode {
	if s.cfg.Section("torrent").HasKey("file_permissions") {
		if fp, err := strconv.ParseUint(s.cfg.Section("torrent").Key("file_permissions").String(), 8, 32); err == nil {
			return fs.FileMode(uint32(fp))
		} else {
			log.Println("Couldn't convert torrent.file_permissions %q: %v", fp, err)
		}
	}

	return fs.FileMode(uint32(0o755))
}

func makeTorrentConfig(s *server) torrent.Config {
	tcfg := torrent.DefaultConfig
	tcfg.RPCEnabled = false
	tcfg.DataDir = s.torrentBaseDir()
	tcfg.DataDirIncludesTorrentID = s.datadirIncludesTorrentID()
	tcfg.Database = s.torrentMetadatadir()
	tcfg.FilePermissions = s.filePermissions()

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
