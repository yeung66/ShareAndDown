package api

import (
	"github.com/gin-gonic/gin"

	"github.com/yeung66/ShareAndDown/utils"
	"strconv"
)

var ResourcePath = "./resources"

var (
	maxSaveMinutes = 20
)

func SetUploadPath(path string) {
	ResourcePath = path
}

func UploadHandler(c *gin.Context) {
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
	savePath := ResourcePath + "/upload/" + token + file.Filename
	err = c.SaveUploadedFile(file, savePath)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	codePath := ResourcePath + "/static/qrcodes/" + token + ".jpg"
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

func DownloadHandler(c *gin.Context) {
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
