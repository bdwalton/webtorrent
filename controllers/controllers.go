package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/anacrolix/torrent/types/infohash"
	"github.com/bdwalton/webtorrent/models"
	"github.com/gin-gonic/gin"
)

func GetTorrents(c *gin.Context) {
	torrents := []*models.Torrent{}
	for _, t := range srv.client.Torrents() {
		torrents = append(torrents, models.FromTorrent(t))
	}

	c.JSON(http.StatusOK, torrents)
}

func AddTorrent(c *gin.Context) {
	var td models.TextData

	if err := c.BindJSON(&td); err != nil {
		c.JSON(http.StatusBadRequest, "")
	}

	uri := td.Data

	if !strings.HasPrefix(uri, "magnet:") {
		c.JSON(http.StatusBadRequest, "")
		return
	}

	log.Printf("Webtorrent: Asked to torrent %q.", uri)
	t, err := srv.client.AddMagnet(uri)
	if err != nil {
		log.Printf("WebTorrent: Error adding magnet uri %q: %v", uri, err)
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	// Return the to client before we do the rest of the setup as
	// that can block for a long time.
	c.JSON(http.StatusOK, models.FromTorrent(t))

	// TODO(bdwalton): This can be problematic as it may return
	// very late or sometimes "never." Put it in a goroutine and
	// treat it as a failure after a timeout? That would still
	// leak whatever resources the torrent client itself is
	// consuming.
	<-t.GotInfo()

	// No fatal errors allowed beyond this point
	srv.trackTorrent(uri, t)
}

func StartTorrent(c *gin.Context) {
	var ti models.Torrent

	if err := c.BindJSON(&ti); err != nil {
		c.JSON(http.StatusBadRequest, "")
	}

	t, ok := srv.client.Torrent(infohash.FromHexString(ti.Hash))
	if !ok {
		c.JSON(http.StatusBadRequest, "")
	}

	t.AllowDataUpload()
	t.AllowDataDownload()
	t.DownloadAll()

	c.JSON(http.StatusOK, "")

}

func PauseTorrent(c *gin.Context) {
	var ti models.Torrent

	if err := c.BindJSON(&ti); err != nil {
		c.JSON(http.StatusBadRequest, "")
	}

	t, ok := srv.client.Torrent(infohash.FromHexString(ti.Hash))
	if !ok {
		c.JSON(http.StatusBadRequest, "")
	}

	t.DisallowDataUpload()
	t.DisallowDataDownload()

	c.JSON(http.StatusOK, "")
}

func DeleteTorrent(c *gin.Context) {
	t, ok := srv.client.Torrent(infohash.FromHexString(c.Param("hash")))
	if !ok {
		c.JSON(http.StatusBadRequest, "")
	}

	srv.dropTorrent(t)

	c.JSON(http.StatusOK, models.FromTorrent(t))
}

func TorrentStatus(c *gin.Context) {
	s := strings.Builder{}
	srv.client.WriteStatus(&s)
	c.JSON(http.StatusOK, models.TextDataFromString(s.String()))
}

func ShowConfig(c *gin.Context) {
	s := strings.Builder{}
	srv.cfg.WriteTo(&s)
	c.JSON(http.StatusOK, models.TextDataFromString(s.String()))
}

func ShutdownTorrentClient() {
	srv.client.Close()
	<-srv.client.Closed()
}
