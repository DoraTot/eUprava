package websocket

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Client struct {
	UserID string
	Conn   *websocket.Conn
}

type Hub struct {
	Clients map[string]*Client
	Mutex   sync.Mutex
}

var WsHub = Hub{
	Clients: make(map[string]*Client),
}

var PendingNotifications = make(map[string][]string)

func SendNotification(userID, message string) {
	WsHub.Mutex.Lock()
	client, ok := WsHub.Clients[userID]
	if ok {
		client.Conn.WriteMessage(websocket.TextMessage, []byte(message))
	} else {
		PendingNotifications[userID] = append(PendingNotifications[userID], message)
	}
	WsHub.Mutex.Unlock()
}
