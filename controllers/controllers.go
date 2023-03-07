package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/bdwalton/webtorrent/models"
	"github.com/cenkalti/rain/torrent"
	"github.com/gin-gonic/gin"
)

func GetTorrents(c *gin.Context) {
	torrents := []*models.BasicTorrentData{}
	for _, t := range srv.client.ListTorrents() {
		torrents = append(torrents, models.BasicTorrentDataFromTorrent(t))
	}

	c.JSON(http.StatusOK, torrents)
}

func AddTorrent(c *gin.Context) {
	var tu models.TorrentURI

	if err := c.BindJSON(&tu); err != nil {
		m := &models.APIError{
			Error:  "Failed to parse request",
			Detail: "Call to AddTorrent() unable to parse input.",
		}
		c.JSON(http.StatusBadRequest, m)
	}

	if !strings.HasPrefix(tu.URI, "magnet:") {
		m := &models.APIError{
			Error:  "Invalid URI.",
			Detail: "Non-magnet URI supplied. We only accept magnet URI.",
		}
		c.JSON(http.StatusBadRequest, m)
	}

	// TODO: Add config knobs for AddTorrentOptions
	t, err := srv.client.AddURI(tu.URI, &torrent.AddTorrentOptions{Stopped: true})
	if err != nil {
		log.Printf("WebTorrent: Error adding URI %q: %v", tu.URI, err)
		m := &models.APIError{
			Error:  "Failed to consume URI",
			Detail: "Call to AddTorrent() unable to add URI..",
		}
		c.JSON(http.StatusInternalServerError, m)
	}

	c.JSON(http.StatusOK, models.BasicTorrentDataFromTorrent(t))
}

func StartTorrent(c *gin.Context) {
	var td models.TorrentID
	if err := c.BindJSON(&td); err != nil {
		log.Printf("WebTorrent: Failed to parse TorrentID: %v", err)
		m := &models.APIError{
			Error:  "Failed to parse request",
			Detail: "Call to StartTorrent() unable to parse input.",
		}
		c.JSON(http.StatusBadRequest, m)
	}

	// if err := srv.startTorrent(td.Data); err != nil {

	c.JSON(http.StatusOK, "")

}

func PauseTorrent(c *gin.Context) {
	var td models.TorrentID

	if err := c.BindJSON(&td); err != nil {
		log.Printf("WebTorrent: Failed to parse TorrentID: %v", err)
		m := &models.APIError{
			Error:  "Failed to parse request",
			Detail: "Call to PauseTorrent() unable to parse input.",
		}
		c.JSON(http.StatusBadRequest, m)
	}

	//	if err := srv.pauseTorrent(td.Data); err != nil {

	c.JSON(http.StatusOK, "")
}

func DeleteTorrent(c *gin.Context) {
	// hash := c.Param("hash")
	// md, ok := srv.torrents[hash]

	c.JSON(http.StatusOK, "")
}

func TorrentDetails(c *gin.Context) {
	// hash := c.Param("hash")

	// if !found {
	// 	m := &models.APIError{
	// 		Error:  "Unknown torrent",
	// 		Detail: fmt.Sprintf("Torrent %q isn't known by the server.", hash),
	// 	}
	// 	c.JSON(http.StatusBadRequest, m)
	// 	return
	// }

	// d, err := srv.torrentDetails(hash)

	c.JSON(http.StatusOK, "")
}

func TorrentStatus(c *gin.Context) {
	s := strings.Builder{}
	c.JSON(http.StatusOK, models.ServerData{Data: s.String()})
}

func ShowConfig(c *gin.Context) {
	s := strings.Builder{}
	srv.cfg.WriteTo(&s)
	c.JSON(http.StatusOK, models.ServerData{Data: s.String()})
}

func ShutdownTorrentClient() error {
	return srv.client.Close()
}
