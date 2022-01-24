package ws

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/yeung66/ShareAndDown/utils/chat"
	"net/http"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	expiredTime = 10 * 60 // 10 minutes
)

func WsHandler(c *gin.Context) {
	wsToken := c.Param("token")
	chatWs, ok := chat.GetWsToken(wsToken)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "wrong websocket token",
		})
		return
	}
	chatWs.ToSend = make(chan chat.WsMessage)

	//chatWs.Expired = make(chan time.Time)

	func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsupgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Failed to set websocket upgrade: %+v", err)
			return
		}

		closed := make(chan bool)
		defer close(closed)

		go func() {
			for {
				_, msg, err := conn.ReadMessage()
				chatWsInRooms, ok1 := chat.GetWsTokenInRoom(chatWs.RoomToken)
				if err != nil || !ok1 {
					closed <- true
					return
				}

				for _, c := range chatWsInRooms {
					if c.WsToken != wsToken {
						c.ToSend <- chat.WsMessage{
							Send:    wsToken,
							Message: string(msg),
						}
					}
				}
			}
		}()

		for {
			select {
			case msg := <-chatWs.ToSend:
				conn.WriteJSON(msg)
			case <-closed:
				return
			}
		}
	}(c.Writer, c.Request)
}
