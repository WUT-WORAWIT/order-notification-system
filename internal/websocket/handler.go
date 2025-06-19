package websocket

import (
	"net/http"
	"order-notification-system/internal/utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	utils.RegisterClient(conn)
	defer utils.UnregisterClient(conn)

	for {
		// Listen for messages from the client if needed
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
