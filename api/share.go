package api

import (
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"

	"github.com/yeung66/ShareAndDown/utils"
	"os"
	"strconv"
)

var route *gin.Engine
var resourcePath string = "./resources"

var (
	maxSaveMinutes       = 20
	port                 = "8000"
	maxBodyBytes   int64 = 25 << 20
)

func InitServer() {
	route = gin.Default()
	route.Use(limits.RequestSizeLimiter(maxBodyBytes))

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	route.MaxMultipartMemory = 20 << 20 // 20Mib

	route.Static("/index", resourcePath+"/html")
	route.Static("/static", resourcePath+"/static")

	sendGroup := route.Group("/share")
	{
		sendGroup.POST("/upload", uploadHandler)
		sendGroup.GET("/download/:token", downloadHandler)
	}

	if p, ok := os.LookupEnv("PORT"); ok {
		port = p
	}
	route.Run("localhost:" + port)
}

func SetUploadPath(path string) {
	resourcePath = path
}

func uploadHandler(c *gin.Context) {
	if !utils.AllowUpload() {
		c.JSON(503, gin.H{
			"status":  "error",
			"message": "too much files uploaded",
		})
		return
	}

	file, err := c.FormFile("file")

	if len(c.Errors) > 0 {
		return
	}

	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if file.Size > (20 << 20) { // 20 Mib
		c.JSON(401, gin.H{
			"status":  "error",
			"message": "too large file with size over 20 Mib",
		})
		return
	}

	saveType := c.DefaultPostForm("saveType", "byCount")
	saveTime := 0
	if saveType == "byTime" {
		saveTimeTemp := c.PostForm("saveTime")
		saveTime, err = strconv.Atoi(saveTimeTemp)
		if err != nil || saveTime <= 0 || saveTime > maxSaveMinutes {
			c.JSON(400, gin.H{
				"status":  "ok",
				"message": "invalid time duration",
			})
			return
		}
	}

	token := utils.TokenGenerator()
	savePath := resourcePath + "/upload/" + token + file.Filename
	err = c.SaveUploadedFile(file, savePath)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	codePath := resourcePath + "/static/qrcodes/" + token + ".jpg"
	codeShow := "http://" + c.Request.Host + "/static/qrcodes/" + token + ".jpg"
	downloadUrl := "http://" + c.Request.Host + "/share/download/" + token
	err = utils.GenQRCode(downloadUrl, codePath)
	if err != nil {
		codePath = ""
		codeShow = ""
	}

	utils.AddFileInfo(file.Filename, savePath, token, codePath, saveTime)
	c.JSON(200, gin.H{
		"status":       "ok",
		"qrcode":       codeShow,
		"fileUrl":      downloadUrl,
		"save_type":    saveType,
		"save_minutes": saveTime,
	})

	if saveTime > 0 {
		utils.DelFileAfter(token, saveTime)
	} else {
		utils.DelFileAfter(token, maxSaveMinutes)
	}

}

func downloadHandler(c *gin.Context) {
	token := c.Param("token")

	fileInfo, err := utils.GetFileInfo(token)
	if err != nil {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+fileInfo.Filename)
	c.Header("Content-Type", "application/octet-stream")

	c.File(fileInfo.FilePath)

	if !fileInfo.SaveFromTime {
		utils.DelFile(token)
	}
}
