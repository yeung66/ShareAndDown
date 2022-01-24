package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yeung66/ShareAndDown/utils/chat"
	"net/http"
)

func CreateChatHandler(c *gin.Context) {
	chatWs := chat.NewChat()
	c.JSON(http.StatusOK, gin.H{
		"chatToken":      chatWs.RoomToken,
		"websocketToken": chatWs.WsToken,
	})
}

func JoinChatHandler(c *gin.Context) {
	roomToken := c.Param("token")
	chatWs, err := chat.JoinChat(roomToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "wrong chat token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chatToken":      roomToken,
		"websocketToken": chatWs.WsToken,
	})

}
