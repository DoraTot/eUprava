package handler

import (
	gws "github.com/gorilla/websocket"
	"log"
	ws "main.go/websocket"
	"net/http"
)

var upgrader = gws.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	client := &ws.Client{
		UserID: userID,
		Conn:   conn,
	}

	ws.WsHub.Mutex.Lock()
	ws.WsHub.Clients[client.UserID] = client

	if msgs, ok := ws.PendingNotifications[client.UserID]; ok {
		for _, m := range msgs {
			client.Conn.WriteMessage(gws.TextMessage, []byte(m))
		}
		delete(ws.PendingNotifications, client.UserID)
	}
	ws.WsHub.Mutex.Unlock()

	go func() {
		defer func() {
			ws.WsHub.Mutex.Lock()
			delete(ws.WsHub.Clients, client.UserID)
			ws.WsHub.Mutex.Unlock()
			conn.Close()
			log.Println("User disconnected:", client.UserID)
		}()

		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}()

	log.Println("User connected:", client.UserID)
}
