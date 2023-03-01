package mappings

import (
	"github.com/bdwalton/webtorrent/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func Init() {
	router = gin.Default()

	// For now, allow all origins. We can tighten this up later.
	router.Use(cors.Default())

	router.LoadHTMLGlob("templates/*.tmpl.html")

	router.GET("/", controllers.Index)
	router.GET("/showconfig", controllers.ShowConfig)
	router.GET("/torrentstatus", controllers.TorrentStatus)

	v1 := router.Group("/v1")
	{
		v1.GET("/torrent", controllers.GetTorrents)
		v1.POST("/torrent", controllers.AddTorrent)
		v1.PUT("/torrent/start", controllers.StartTorrent)
		v1.PUT("/torrent/pause", controllers.PauseTorrent)
		v1.DELETE("/torrent/:hash", controllers.DeleteTorrent)
	}
}

func GetRouter() *gin.Engine {
	return router
}
