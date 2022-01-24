package chat

import (
	"fmt"
	"github.com/yeung66/ShareAndDown/utils"
	"time"
)

type ChatWebSocket struct {
	RoomToken, WsToken string
	Expired            chan time.Time
	ToSend             chan WsMessage
}

type WsMessage struct {
	Send    string `json:"send"`
	Message string `json:"message"`
}

var chatsRoom = map[string][]*ChatWebSocket{}
var chats2Room = map[string]*ChatWebSocket{}

func NewChat() *ChatWebSocket {
	token := utils.TokenGenerator()
	for _, ok := chatsRoom[token]; ok; {
		token = utils.TokenGenerator()
	}

	wsToken := utils.TokenGenerator()
	chatWs := &ChatWebSocket{
		RoomToken: token,
		WsToken:   wsToken,
	}

	chatsRoom[token] = []*ChatWebSocket{chatWs}
	chats2Room[wsToken] = chatWs

	return chatWs
}

func JoinChat(roomToken string) (*ChatWebSocket, error) {
	room, ok := chatsRoom[roomToken]
	if !ok {
		return &ChatWebSocket{}, fmt.Errorf("no such chat room")
	}

	wsToken := utils.TokenGenerator()
	chatWs := &ChatWebSocket{
		RoomToken: roomToken,
		WsToken:   wsToken,
	}

	chatsRoom[roomToken] = append(room, chatWs)

	chats2Room[wsToken] = chatWs
	return chatWs, nil
}

func GetWsToken(wsToken string) (*ChatWebSocket, bool) {
	chatWs, ok := chats2Room[wsToken]
	return chatWs, ok
}

func GetWsTokenInRoom(roomToken string) ([]*ChatWebSocket, bool) {
	chatsInRoom, ok := chatsRoom[roomToken]
	return chatsInRoom, ok
}
