package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/types/infohash"
	"github.com/bdwalton/webtorrent/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
)

type server struct {
	client *torrent.Client
	cfg    *ini.File
}

var srv *server

func Init(cfg *ini.File) error {
	tcfg := torrent.NewDefaultClientConfig()
	tcfg.DataDir = cfg.Section("torrent").Key("datadir").String()

	c, err := torrent.NewClient(tcfg)
	if err != nil {
		return fmt.Errorf("Error establishing torrent client: %v\n", err)
	}

	srv = &server{
		client: c,
		cfg:    cfg,
	}

	return nil
}

func torrentInfoFromTorrent(t *torrent.Torrent) *models.TorrentInfo {
	return &models.TorrentInfo{
		Name: t.Name(),
		Hash: t.InfoHash().HexString(),
	}
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl.html", "")
}

func GetTorrents(c *gin.Context) {
	torrents := []*models.TorrentInfo{}
	for _, t := range srv.client.Torrents() {
		torrents = append(torrents, torrentInfoFromTorrent(t))
	}

	c.JSON(http.StatusOK, torrents)
}

func AddTorrent(c *gin.Context) {
	var taddr models.TorrentAddress

	if err := c.BindJSON(&taddr); err != nil {
		c.JSON(http.StatusBadRequest, "")
	}

	if !strings.HasPrefix(taddr.URI, "magnet:") {
		c.JSON(http.StatusBadRequest, "")
	}

	t, err := srv.client.AddMagnet(taddr.URI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
	}

	<-t.GotInfo()
	// t.DownloadAll()
	c.JSON(http.StatusOK, torrentInfoFromTorrent(t))
}

func DeleteTorrent(c *gin.Context) {
	t, ok := srv.client.Torrent(infohash.FromHexString(c.Param("hash")))
	if !ok {
		c.JSON(http.StatusBadRequest, "")
	}

	t.Drop()
	c.JSON(http.StatusOK, "")
}

// func ToggleTorrent(c *gin.Context) {
// 	c.JSON(http.StatusOK, "")
// }

func TorrentStatus(c *gin.Context) {
	s := strings.Builder{}
	srv.client.WriteStatus(&s)
	c.HTML(http.StatusOK, "torrentstatus.tmpl.html", s.String())
}

func ShowConfig(c *gin.Context) {
	s := strings.Builder{}
	srv.cfg.WriteTo(&s)
	c.HTML(http.StatusOK, "showconfig.tmpl.html", s.String())
}

func ShutdownTorrentClient() {
	srv.client.Close()
	<-srv.client.Closed()
}