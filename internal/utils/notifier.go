package utils

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
)

// NotifyNewOrder broadcasts a new order notification to all connected WebSocket clients.
func NotifyNewOrder(orderID string, itemCode string, item string, quantity int) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	message := map[string]interface{}{
		"orderID":  orderID,
		"itemCode": itemCode,
		"item":     item,
		"quantity": quantity,
	}

	for client := range clients {
		err := client.WriteJSON(message)
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}

// RegisterClient adds a new WebSocket client to the list of connected clients.
func RegisterClient(client *websocket.Conn) {
	clientsMu.Lock()
	clients[client] = true
	clientsMu.Unlock()
}

// UnregisterClient removes a WebSocket client from the list of connected clients.
func UnregisterClient(client *websocket.Conn) {
	clientsMu.Lock()
	delete(clients, client)
	clientsMu.Unlock()
}
