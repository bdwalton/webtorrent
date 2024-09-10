package mappings

import (
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/bdwalton/webtorrent/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/ini.v1"
)

var router *gin.Engine

func Init(cfg *ini.File, basePath string, staticFiles fs.FS) {
	gin.SetMode(ginMode(cfg))
	router = gin.Default()

	// For now, allow all origins. We can tighten this up later.
	router.Use(cors.Default())

	defPage := filepath.Join(basePath, "webtorrent.html")
	hfs := http.FS(staticFiles)
	router.StaticFileFS("/", defPage, hfs)
	router.StaticFileFS("/webtorrent.html", defPage, hfs)

	// Maybe provide an embed.FS for this later, but for now, we
	// can serve them from the filesystem.
	router.NoRoute(func(c *gin.Context) {
		page := basePath + path.Join(c.Request.RequestURI)
		// If the user force refreshes their browser while one
		// of the virtual angular SPA endpoints is in the
		// location bar, the browser will request that path
		// from us. We cover that up by defaulting back to
		// serving the underlying angular html page.
		if _, err := os.Stat(page); errors.Is(err, os.ErrNotExist) {
			c.FileFromFS(defPage, hfs)
		}
		c.FileFromFS(page, hfs)
	})

	router.StaticFileFS("/signin", filepath.Join(basePath, "assets/signin.html"), hfs)
	router.GET("/signout", controllers.SignoutHandler)
	router.GET("/auth/:provider", controllers.SignInWithProvider)
	router.GET("/auth/:provider/callback", controllers.CallBackHandler)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := router.Group("/v1")
	{
		// Torrent interaction calls
		v1.GET("/torrent", controllers.GetTorrents)
		v1.POST("/torrent", controllers.AddTorrent)
		v1.PUT("/torrent/start", controllers.StartTorrent)
		v1.PUT("/torrent/pause", controllers.PauseTorrent)
		v1.DELETE("/torrent/:id", controllers.DeleteTorrent)

		// Torrent details interaction calls
		v1.GET("/torrentdetails/:id", controllers.TorrentDetails)

		// Server health diagnostic calls
		v1.GET("/showconfig", controllers.ShowConfig)
		v1.GET("/torrentstatus", controllers.TorrentClientStatus)
	}
}

func GetRouter() *gin.Engine {
	return router
}

func ginMode(cfg *ini.File) string {
	if cfg.Section("server").HasKey("gin_mode") {
		return cfg.Section("server").Key("gin_mode").String()
	}

	return "release"
}
