package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yeung66/ShareAndDown/utils"
	"strconv"
)

var route *gin.Engine
var resourcePath string = "./resources"

func InitServer() {
	route = gin.Default()

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	route.MaxMultipartMemory = 20 << 20 // 20Mib

	route.Static("/index", resourcePath+"/html")
	route.Static("/static", resourcePath+"/static")

	sendGroup := route.Group("/share")
	{
		sendGroup.POST("/upload", uploadHandler)
		sendGroup.GET("/download/:token", downloadHandler)
	}

	route.Run(":8000")
}

func SetUploadPath(path string) {
	resourcePath = path
}

func uploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
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

	codePath := resourcePath + "/static/qrcodes/" + token + file.Filename + ".jpg"
	codeShow := "http://" + c.Request.Host + "/static/qrcodes/" + token + file.Filename + ".jpg"
	downloadUrl := "http://" + c.Request.Host + "/share/download/" + token
	err = utils.GenQRCode(downloadUrl, codePath)
	if err != nil {
		codePath = ""
		codeShow = ""
	}

	saveType := c.DefaultPostForm("saveType", "byCount")
	saveTime := 0
	if saveType == "byTime" {
		saveTimeTemp := c.PostForm("saveTime")
		saveTime, err = strconv.Atoi(saveTimeTemp)
		if err != nil || saveTime <= 0 || saveTime > 20 {
			c.JSON(401, gin.H{
				"status":  "ok",
				"message": "invalid time duration",
			})
			return
		}
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
