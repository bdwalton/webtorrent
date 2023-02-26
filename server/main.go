// package server implements the webtorrent http server.
package server

import (
	"context"
	"embed"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/anacrolix/torrent"
	"gopkg.in/ini.v1"
)

//go:embed templates/*.tmpl.html
var templateFS embed.FS

type torrentServer struct {
	c     *torrent.Client
	cfg   *ini.File
	tmpls *template.Template
}

func newTorrentServer(c *torrent.Client, cfg *ini.File) *torrentServer {
	return &torrentServer{
		c:   c,
		cfg: cfg,
		// This may panic, but only during startup, so should be fine.
		tmpls: template.Must(template.ParseFS(templateFS, "templates/*.tmpl.html")),
	}
}

func (ts *torrentServer) rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(ts.tmpls)
	if err := ts.tmpls.ExecuteTemplate(w, "root.tmpl.html", nil); err != nil {
		log.Println("TorrentServer: Failed to execute template:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to complete response. See server logs for details."))
	}
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
