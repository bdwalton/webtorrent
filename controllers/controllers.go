package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/anacrolix/torrent/types/infohash"
	"github.com/bdwalton/webtorrent/models"
	"github.com/gin-gonic/gin"
)

func GetTorrents(c *gin.Context) {
	torrents := []*models.Torrent{}
	for _, md := range srv.torrents {
		torrents = append(torrents, models.FromTorrentData(md))
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
		m := &models.APIError{
			Error:  "Failed adding URI",
			Detail: "The torrent library was unable or unwilling to consume your URI. See error logs for details",
		}
		c.JSON(http.StatusInternalServerError, m)
		return
	}

	md := srv.registerTorrent(uri, t)

	// Return the to client before we do the rest of the setup as
	// that can block for a long time.
	c.JSON(http.StatusOK, models.FromTorrentData(md))

	// TODO(bdwalton): This can be problematic as it may return
	// very late or sometimes "never." Put it in a goroutine and
	// treat it as a failure after a timeout? That would still
	// leak whatever resources the torrent client itself is
	// consuming.
	<-t.GotInfo()
}

func StartTorrent(c *gin.Context) {
	var ti models.Torrent

	if err := c.BindJSON(&ti); err != nil {
		m := &models.APIError{
			Error:  "Failed to parse request",
			Detail: "Call to StartTorrent() unable to parse input.",
		}
		c.JSON(http.StatusBadRequest, m)
	}

	if err := srv.startTorrent(ti.Hash); err != nil {
		log.Printf("WebTorrent: Failed to start torrent: %v", err)
		m := &models.APIError{
			Error:  "Failed to start torrent.",
			Detail: "Call to StartTorrent() failed. See error logs for details.",
		}
		c.JSON(http.StatusBadRequest, m)
		return
	}

	c.JSON(http.StatusOK, "")

}

func PauseTorrent(c *gin.Context) {
	var ti models.Torrent

	if err := c.BindJSON(&ti); err != nil {
		m := &models.APIError{
			Error:  "Failed to parse request",
			Detail: "Call to PauseTorrent() unable to parse input.",
		}
		c.JSON(http.StatusBadRequest, m)
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
	hash := c.Param("hash")
	md, ok := srv.torrents[hash]
	if !ok {
		m := &models.APIError{
			Error:  "Unknown torrent",
			Detail: fmt.Sprintf("Torrent %q isn't known by the server.", hash),
		}
		c.JSON(http.StatusBadRequest, m)
		return
	}

	if err := srv.dropTorrent(hash); err != nil {
		log.Printf("WebTorrent: Error dropping torrent: %v", err)
		m := &models.APIError{
			Error:  "Failed to delete torrent",
			Detail: "Cleaning up the torrent failed, see error logs for details.",
		}
		c.JSON(http.StatusInternalServerError, m)
		return
	}

	c.JSON(http.StatusOK, models.FromTorrentData(md))
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
