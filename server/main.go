// package server implements the webtorrent http server.
package server

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/types/infohash"
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

func generic500(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Unable to complete reponse. See server logs for details."))
}

type rootData struct {
	Torrents []*torrent.Torrent
}

func (ts *torrentServer) rootHandler(w http.ResponseWriter, r *http.Request) {
	data := rootData{
		Torrents: ts.c.Torrents(),
	}
	if err := ts.tmpls.ExecuteTemplate(w, "root.tmpl.html", data); err != nil {
		log.Println("TorrentServer: Failed to execute template:", err)
		generic500(w)
	}
}

func (ts *torrentServer) addTorrentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		r.ParseForm()
		muri := r.FormValue("torrenturi")
		if strings.HasPrefix(muri, "magnet:") {
			t, err := ts.c.AddMagnet(muri)
			if err != nil {
				log.Println("TorrentServer: Couldn't add magnet URI:", err)
				generic500(w)
			}
			go func(t *torrent.Torrent) {
				<-t.GotInfo()
				t.DownloadAll()
			}(t)
		} else {
			log.Println("TorrentServer: Invalid URI:", muri)
		}

		http.Redirect(w, r, r.URL.Host+"/", http.StatusTemporaryRedirect)
	default:
		log.Println("TorrentServer: Received non-POST request to", r.RequestURI)
		generic500(w)
	}
}

func (ts *torrentServer) pauseTorrentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		r.ParseForm()
		h := r.FormValue("torrenthash")
		log.Printf("TorrentServer: Asked to pause torrent with hash: %q", h)
		t, ok := ts.c.Torrent(infohash.FromHexString(h))
		if !ok {
			log.Printf("TorrentServer: Invalid torrent hash %q", h)
			generic500(w)
			return
		}

		t.DisallowDataUpload()
		t.DisallowDataDownload()
		w.Write([]byte("Paused " + h))
		return
	default:
		log.Println("TorrentServer: Received non-POST request to", r.RequestURI)
		generic500(w)
	}
}

func (ts *torrentServer) resumeTorrentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		r.ParseForm()
		h := r.FormValue("torrenthash")
		log.Printf("TorrentServer: Asked to resume torrent with hash %q", h)
		t, ok := ts.c.Torrent(infohash.FromHexString(h))
		if !ok {
			log.Printf("TorrentServer: Invalid torrent hash %q", h)
			generic500(w)
			return
		}

		t.AllowDataUpload()
		t.AllowDataDownload()
		w.Write([]byte("Resumed " + h))
		return
	default:
		log.Println("TorrentServer: Received non-POST request to", r.RequestURI)
		generic500(w)
	}
}

func (ts *torrentServer) stopTorrentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		r.ParseForm()
		h := r.FormValue("torrenthash")
		log.Printf("TorrentServer: Asked to stop torrent with hash %q", h)
		t, ok := ts.c.Torrent(infohash.FromHexString(h))
		if !ok {
			log.Printf("TorrentServer: Invalid torrent hash %q", h)
			generic500(w)
			return
		}

		t.Drop()
		w.Write([]byte("Stopped " + h))
		return
	default:
		log.Println("TorrentServer: Received non-POST request to", r.RequestURI)
		generic500(w)
	}
}

func (ts *torrentServer) quitHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("TorrentServer: /quitquitquit requested, so shutting down...")
	fmt.Fprintln(w, "<p>Goodbye...</p>")

	ts.c.Close()
}

func (ts *torrentServer) clientStatusHandler(w http.ResponseWriter, r *http.Request) {
	s := strings.Builder{}
	ts.c.WriteStatus(&s)

	if err := ts.tmpls.ExecuteTemplate(w, "torrentstatus.tmpl.html", s.String()); err != nil {
		log.Println("TorrentServer: Failed to execute /torrentstatus template:", err)
		generic500(w)
	}
}

func (ts *torrentServer) serve(ctx context.Context, s *http.Server) error {
	http.HandleFunc("/", ts.rootHandler)
	http.HandleFunc("/addtorrent", ts.addTorrentHandler)
	http.HandleFunc("/pausetorrent", ts.pauseTorrentHandler)
	http.HandleFunc("/resumetorrent", ts.resumeTorrentHandler)
	http.HandleFunc("/stoptorrent", ts.stopTorrentHandler)
	http.HandleFunc("/clientstatus", ts.clientStatusHandler)
	http.HandleFunc("/quitquitquit", ts.quitHandler)

	log.Printf("TorrentServer: Listening on %q", s.Addr)
	return s.ListenAndServe()
}

func ListenAndServe(ctx context.Context, s *http.Server, c *torrent.Client, cfg *ini.File) error {
	ts := newTorrentServer(c, cfg)
	return ts.serve(ctx, s)
}
