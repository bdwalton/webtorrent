package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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

	md, err := srv.registerTorrent(uri)
	if err != nil {
		log.Printf("WebTorrent: StartTorrent() call to registerTorrent() failed: %v", err)
		m := &models.APIError{
			Error:  "Failed to register URI.",
			Detail: "Unable to register the URI in the server. See error logs for details.",
		}
		c.JSON(http.StatusInternalServerError, m)
	}

	// Return the to client before we do the rest of the setup as
	// that can block for a long time.
	c.JSON(http.StatusOK, models.FromTorrentData(md))

	// TODO(bdwalton): This can be problematic as it may return
	// very late or sometimes "never." Put it in a goroutine and
	// treat it as a failure after a timeout? That would still
	// leak whatever resources the torrent client itself is
	// consuming.
	<-md.T.GotInfo()
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

	if err := srv.pauseTorrent(ti.Hash); err != nil {
		log.Printf("WebTorrent: Failed to pause torrent: %v", err)
		m := &models.APIError{
			Error:  "Failed to pause torrent",
			Detail: "Call to PauseTorrent() unable to complete request. See error logs for details.",
		}
		c.JSON(http.StatusInternalServerError, m)
	}

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
