package web

import (
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/yeung66/ShareAndDown/web/api"
	"github.com/yeung66/ShareAndDown/web/ws"
	"os"
)

var route *gin.Engine

var (
	port               = "8000"
	maxBodyBytes int64 = 25 << 20
)

func InitServer() {
	var resourcePath = api.ResourcePath

	route = gin.Default()
	route.Use(limits.RequestSizeLimiter(maxBodyBytes))

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	route.MaxMultipartMemory = 20 << 20 // 20Mib

	route.Static("/index", resourcePath+"/html")
	route.Static("/static", resourcePath+"/static")

	sendGroup := route.Group("/share")
	{
		sendGroup.POST("/upload", api.UploadHandler)
		sendGroup.GET("/download/:token", api.DownloadHandler)
	}

	route.GET("/chat", api.CreateChatHandler)
	route.GET("/chat/:token", api.JoinChatHandler)
	route.GET("/ws/chat/:token", ws.WsHandler)

	if p, ok := os.LookupEnv("PORT"); ok {
		port = p
	}
	route.Run("localhost:" + port)
}
