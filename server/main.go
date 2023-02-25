// package server implements the webtorrent http server.
package server

import (
	"context"
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"

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

func (ts *torrentServer) quitHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("TorrentServer: /quitquitquit requested, so shutting down...")
	fmt.Fprintln(w, "<p>Goodbye...</p>")

	ts.c.Close()
}

func (ts *torrentServer) torrentStatus(w http.ResponseWriter, r *http.Request) {
	s := strings.Builder{}
	ts.c.WriteStatus(&s)
	out := strings.Replace(html.EscapeString(s.String()), "\n", "<br>\n", -1)
	fmt.Fprintln(w, "<h1>Torrent Client Status</h1>")
	fmt.Fprintln(w, "<tt>", out, "</tt>")
}

func (ts *torrentServer) serve(ctx context.Context, s *http.Server) error {
	http.HandleFunc("/", ts.rootHandler)
	http.HandleFunc("/clientstatus", ts.torrentStatus)
	http.HandleFunc("/quitquitquit", ts.quitHandler)

	log.Printf("TorrentServer: Listening on %q", s.Addr)
	return s.ListenAndServe()
}

func ListenAndServe(ctx context.Context, s *http.Server, c *torrent.Client, cfg *ini.File) error {
	ts := newTorrentServer(c, cfg)
	return ts.serve(ctx, s)
}
