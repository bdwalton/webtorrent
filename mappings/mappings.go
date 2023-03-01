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
		v1.GET("/torrents", controllers.GetTorrents)
		v1.POST("/torrents", controllers.AddTorrent)
		v1.DELETE("/torrents/:hash", controllers.DeleteTorrent)
		// v1.PUT("/torrents/:hash", controllers.ToggleDownload)
	}
}

func GetRouter() *gin.Engine {
	return router
}
