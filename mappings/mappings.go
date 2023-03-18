package mappings

import (
	"path"
	"path/filepath"

	"github.com/bdwalton/webtorrent/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var router *gin.Engine

func Init(ginMode string) {
	gin.SetMode(ginMode)
	router = gin.Default()

	// For now, allow all origins. We can tighten this up later.
	router.Use(cors.Default())

	// Maybe provide an embed.FS for this later, but for now, we
	// can serve them from the filesystem.
	router.NoRoute(func(c *gin.Context) {
		dir, file := path.Split(c.Request.RequestURI)
		ext := filepath.Ext(file)
		if file == "" || ext == "" {
			c.File("./ui/dist/ui/index.html")
		} else {
			c.File("./ui/dist/ui/" + path.Join(dir, file))
		}
	})

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
		v1.GET("/torrentstatus", controllers.TorrentStatus)
	}
}

func GetRouter() *gin.Engine {
	return router
}
