package controllers

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/bdwalton/webtorrent/models"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/protobuf/proto"
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
		func() float64 { return float64(len(srv.client.Torrents())) },
	)

	prometheus.MustRegister(loadedTorrents)
}

type server struct {
	client   *torrent.Client
	cfg      *ini.File
	mtx      sync.Mutex
	torrents map[string]*models.BasicMetaData // Torrent Hash to our additional info
}

// datadir returns the torrent.datadir key from the config as a
// string. This is a small helper because we reference this
// frequently.
func (s *server) datadir() string {
	return s.cfg.Section("torrent").Key("datadir").String()
}

func (s *server) trackTorrent(uri string, t *torrent.Torrent) {
	md := &models.BasicMetaData{
		URI:     uri,
		Running: false,
		T:       t,
	}

	h := t.InfoHash().HexString()
	s.torrents[h] = md

	log.Printf("WebTorrent: Tracking %s", t.String())

	// Errors from this are non-fatal, so nothing returned.
	s.writeMetaInfo(h)
}

// writeMetaInfo stores the torrent's metadata to disk so it can be
// restored later.
func (s *server) writeMetaInfo(hash string) {
	td := s.torrents[hash]
	// Should be true before this is called, but doesn't hurt.
	<-td.T.GotInfo()

	data := &models.WebTorrentMetadata{
		Hash:    proto.String(hash),
		Uri:     proto.String(td.URI),
		Running: proto.Bool(td.Running),
	}

	var buf bytes.Buffer
	if err := bencode.NewEncoder(&buf).Encode(td.T.Metainfo()); err != nil {
		log.Printf("WebTorrent: File to write torrent metadata for %q: %v", hash, err)
		return
	}
	data.TorrentInfo = buf.Bytes()

	bin, err := proto.Marshal(data)
	if err != nil {
		log.Printf("WebTorrent: Error serializing proto prior to storing it: %v", err)
		return
	}

	// Protect the create and write operations.
	s.mtx.Lock()
	defer s.mtx.Unlock()

	p := filepath.Join(s.datadir(), hash+suffix)
	log.Printf("WebTorrent: Storing metadata for %q in %q.", hash, p)

	// Assume umask gets applied here.
	if err := os.WriteFile(p, bin, 0666); err != nil {
		log.Printf("WebTorrent: Error persisting data for %q: %v", hash, err)
	}
}

func (s *server) dropTorrent(t *torrent.Torrent) {
	t.Drop()
	f := filepath.Join(s.datadir(), t.InfoHash().HexString()+suffix)
	log.Printf("WebTorrent: Removing metainfo file %q.", f)
	if err := os.Remove(f); err != nil {
		// Not fatal, so log and carry on.
		log.Printf("WebTorrent: Error removing metainfo file %q: %v", f, err)
	}

}

// loadMetaInfoFiles will find and load all metadata files in
// srv.datadir() that were previously persisted. It should only be
// called at startup. No locking is done, although it should be safe
// to call this in parallel regardless because the torrent client
// locks internally.
func (s *server) loadMetaInfoFiles() {
	glob := filepath.Join(s.datadir(), "*"+suffix)
	files, err := filepath.Glob(glob)
	// For now, we log and carry on. We should export some metrics
	// for this condition and maybe allow configuration options to
	// control the behaviour.
	if err != nil {
		log.Printf("WebTorrent: Failed to find meta files in %q: %v", glob, err)
		return
	}

	for _, f := range files {
		var td models.WebTorrentMetadata

		log.Printf("WebTorrent: Loading metainfo from %q.", f)

		bin, err := os.ReadFile(f)
		if err != nil {
			log.Printf("WebTorrent: Failed to load metainfo file %q: %v", f, err)
			continue
		}

		if err := proto.Unmarshal(bin, &td); err != nil {
			log.Printf("WebTorrent: Failed to unmarshall data from file %q: %v", f, err)
			continue
		}

		mi, err := metainfo.Load(bytes.NewBuffer(td.GetTorrentInfo()))
		if err != nil {
			log.Printf("WebTorrent: Error loading metainfo from proto: %v", err)
			continue
		}

		t, err := srv.client.AddTorrent(mi)
		if err != nil {
			log.Printf("WebTorrent: Error instantiating metainfo: %v", err)
			continue
		}

		s.torrents[td.GetHash()] = &models.BasicMetaData{
			URI:     td.GetUri(),
			Running: td.GetRunning(),
			T:       t,
		}

		log.Printf("WebTorrent: Loaded %s from metainfo in %q.", t.String(), f)
	}
}

// Init will create the global srv object and populate it with a
// Torrent client. It also handles initializing pre-saved torrents
// from storage.
func Init(cfg *ini.File) error {
	srv = &server{
		cfg:      cfg,
		torrents: make(map[string]*models.BasicMetaData),
	}

	tcfg := torrent.NewDefaultClientConfig()
	tcfg.DataDir = srv.datadir()

	c, err := torrent.NewClient(tcfg)
	if err != nil {
		return fmt.Errorf("Error establishing torrent client: %v\n", err)
	}

	srv.client = c

	// We consider the torrent loading an optional so errors are
	// swallowed (with logging), but not considered fatal.
	srv.loadMetaInfoFiles()

	registerPrometheus()

	return nil
}
