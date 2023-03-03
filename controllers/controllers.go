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
	var ti models.Torrent

	if err := c.BindJSON(&ti); err != nil {
		c.JSON(http.StatusBadRequest, "")
	}

	if !strings.HasPrefix(ti.URI, "magnet:") {
		c.JSON(http.StatusBadRequest, "")
	}

	log.Printf("Webtorrent: Asked to torrent %q.", ti.URI)
	t, err := srv.client.AddMagnet(ti.URI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
	}

	<-t.GotInfo()
	c.JSON(http.StatusOK, models.FromTorrent(t))

	// Errors from this are non-fatal, so nothing returned.
	srv.writeMetaInfo(t)
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
