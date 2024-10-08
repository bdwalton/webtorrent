package controllers

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/bdwalton/webtorrent/models"
	"github.com/cenkalti/rain/torrent"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func requireSignin(c *gin.Context) {
	username, err := gothic.GetFromSession("username", c.Request)
	if err != nil {
		c.JSON(http.StatusForbidden, "/signin")
		return
	}

	if _, ok := allowedUsers[username]; !ok {
		c.JSON(http.StatusForbidden, "/signin")
		return
	}
}

func GetTorrents(c *gin.Context) {
	requireSignin(c)
	torrents := []models.BasicTorrentData{}
	for _, t := range srv.client.ListTorrents() {
		torrents = append(torrents, models.BasicTorrentDataFromTorrent(t))
	}

	c.JSON(http.StatusOK, torrents)
}

func AddTorrent(c *gin.Context) {
	requireSignin(c)

	var tu models.TorrentURI

	if err := c.BindJSON(&tu); err != nil {
		m := &models.APIError{
			Details: "Call to AddTorrent() unable to parse input.",
		}
		c.JSON(http.StatusBadRequest, m)
	}

	ato := &torrent.AddTorrentOptions{
		Stopped:           srv.stopOnAdd(),
		StopAfterDownload: srv.stopAfterDownload(),
		StopAfterMetadata: srv.stopAfterMetadata(),
	}

	t, err := srv.client.AddURI(tu.URI, ato)
	if err != nil {
		log.Printf("WebTorrent: Error adding URI %q: %v", tu.URI, err)
		m := &models.APIError{
			Details: "Call to AddTorrent() unable to add URI..",
		}
		c.JSON(http.StatusInternalServerError, m)
	}

	go func(wg *sync.WaitGroup) {
		wg.Add(1)
		srv.watchTorrent(t)
		wg.Done()
	}(&srv.wg)

	c.JSON(http.StatusOK, models.BasicTorrentDataFromTorrent(t))
}

func StartTorrent(c *gin.Context) {
	requireSignin(c)

	var td models.TorrentID
	if err := c.BindJSON(&td); err != nil {
		log.Printf("WebTorrent: Failed to parse TorrentID: %v", err)
		m := &models.APIError{
			Details: "Call to StartTorrent() unable to parse input.",
		}
		c.JSON(http.StatusBadRequest, m)
		return
	}

	t := srv.client.GetTorrent(td.ID)
	if t == nil {
		log.Printf("WebTorrent: Unknown Torrent ID %q: %v", td.ID)
		m := &models.APIError{
			Details: "Unknown Torrent ID.",
		}
		c.JSON(http.StatusBadRequest, m)
		return
	}

	if err := t.Start(); err != nil {
		log.Printf("WebTorrent: Failed to start Torrent %q: %v", td.ID, err)
		m := &models.APIError{
			Details: fmt.Sprintf("Failed to start Torrent ID %q.", td.ID),
		}
		c.JSON(http.StatusInternalServerError, m)
		return
	}

	go srv.watchTorrent(t)

	c.JSON(http.StatusOK, models.BasicTorrentDataFromTorrent(t))
}

func PauseTorrent(c *gin.Context) {
	requireSignin(c)

	var td models.TorrentID
	if err := c.BindJSON(&td); err != nil {
		log.Printf("WebTorrent: Failed to parse TorrentID: %v", err)
		m := &models.APIError{
			Details: "Call to PauseTorrent() unable to parse input.",
		}
		c.JSON(http.StatusBadRequest, m)
		return
	}

	t := srv.client.GetTorrent(td.ID)
	if t == nil {
		log.Printf("WebTorrent: Unknown Torrent ID %q: %v", td.ID)
		m := &models.APIError{
			Details: fmt.Sprintf("Unknown Torrent ID %q.", td.ID),
		}
		c.JSON(http.StatusBadRequest, m)
		return
	}

	if err := t.Stop(); err != nil {
		log.Printf("WebTorrent: Failed to pause Torrent %q: %v", td.ID, err)
		m := &models.APIError{
			Details: fmt.Sprintf("Failed to pause Torrent ID %q.", td.ID),
		}
		c.JSON(http.StatusInternalServerError, m)
		return
	}

	c.JSON(http.StatusOK, models.BasicTorrentDataFromTorrent(t))
}

func DeleteTorrent(c *gin.Context) {
	requireSignin(c)

	tid := c.Param("id")
	t := srv.client.GetTorrent(tid)
	if t == nil {
		log.Printf("WebTorrent: Unknown Torrent ID %q", tid)
		m := &models.APIError{
			Details: fmt.Sprintf("Unknown Torrent ID %q.", tid),
		}
		c.JSON(http.StatusBadRequest, m)
		return
	}

	if err := srv.client.RemoveTorrent(tid); err != nil {
		log.Printf("WebTorrent: Failed to remove Torrent %q: %v", tid, err)
		m := &models.APIError{
			Details: fmt.Sprintf("Failed to remove Torrent ID %q.", tid),
		}
		c.JSON(http.StatusInternalServerError, m)
		return
	}

	c.JSON(http.StatusOK, models.BasicTorrentDataFromTorrent(t))
}

func TorrentDetails(c *gin.Context) {
	requireSignin(c)

	tid := c.Param("id")

	t := srv.client.GetTorrent(tid)
	if t == nil {
		log.Printf("WebTorrent: Unknown Torrent ID %q", tid)
		m := &models.APIError{
			Details: fmt.Sprintf("Unknown Torrent ID %q.", tid),
		}
		c.JSON(http.StatusBadRequest, m)
		return
	}

	c.JSON(http.StatusOK, models.TorrentDataFromTorrent(t))
}

func TorrentClientStatus(c *gin.Context) {
	requireSignin(c)

	c.JSON(http.StatusOK, models.WrapSessionStats(srv.client.Stats()))
}
