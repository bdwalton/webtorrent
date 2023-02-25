// package server implements the webtorrent http server.
package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/anacrolix/torrent"
	"gopkg.in/ini.v1"
)

type torrentServer struct {
	c   *torrent.Client
	cfg *ini.File
}

func newTorrentServer(c *torrent.Client, cfg *ini.File) *torrentServer {
	return &torrentServer{
		c:   c,
		cfg: cfg,
	}
}

func (ts *torrentServer) rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello, world!</h1>")
}

func (ts *torrentServer) serve(ctx context.Context) error {
	http.HandleFunc("/", ts.rootHandler)
	return http.ListenAndServe(":"+ts.cfg.Section("server").Key("port").String(), nil)
}

func ListenAndServe(ctx context.Context, c *torrent.Client, cfg *ini.File) error {
	ts := newTorrentServer(c, cfg)
	return ts.serve(ctx)
}
