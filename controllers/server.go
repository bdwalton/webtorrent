package controllers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
	"gopkg.in/ini.v1"
)

const (
	suffix = ".meta"
)

type server struct {
	client *torrent.Client
	cfg    *ini.File
	mtx    sync.Mutex
}

// srv is the package-level object we can reference to find the
// Torrent client and global config.
var srv *server

// datadir returns the torrent.datadir key from the config as a
// string. This is a small helper because we reference this
// frequently.
func (s *server) datadir() string {
	return s.cfg.Section("torrent").Key("datadir").String()
}

// writeMetaInfo stores the torrent's metadata to disk so it can be
// restored later.
func (s *server) writeMetaInfo(t *torrent.Torrent) {
	// Should be true before this is called, but doesn't hurt.
	<-t.GotInfo()

	mi := t.Metainfo()
	h := t.InfoHash().HexString()
	p := filepath.Join(s.datadir(), h+suffix)

	log.Printf("WebTorrent: Storing metadata for %q in %q.", h, p)

	// Protect the create and write operations.
	s.mtx.Lock()
	defer s.mtx.Unlock()

	f, err := os.Create(p)
	if err != nil {
		log.Printf("WebTorrent: Failed to create torrent metadata file for %q: %v", t.InfoHash().HexString(), err)
		return
	}
	defer f.Close()

	if err := bencode.NewEncoder(f).Encode(mi); err != nil {
		log.Printf("WebTorrent: File to write torrent metadata for %q: %v", t.InfoHash().HexString(), err)
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
		mi, err := metainfo.LoadFromFile(f)
		// For now, we log and carry on. We should export some
		// metrics for this condition and maybe allow
		// configuration options to control the behaviour.
		if err != nil {
			log.Printf("WebTorrent: Failed to load metainfo file %q: %v", f, err)
			continue
		}

		t, err := srv.client.AddTorrent(mi)
		if err != nil {
			log.Printf("WebTorrent: Error loading metainfo: %v", err)
		}

		log.Printf("WebTorrent: Loaded %s from metainfo.", t.String())
	}
}

// Init will create the global srv object and populate it with a
// Torrent client. It also handles initializing pre-saved torrents
// from storage.
func Init(cfg *ini.File) error {
	srv = &server{cfg: cfg}

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

	return nil
}
