package chat_ws

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients map[string][]*Client
	mu      *sync.RWMutex
}

type Client struct {
	RideID string
	Conn   *websocket.Conn
	UserID string
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string][]*Client),
		mu:      &sync.RWMutex{},
	}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c.RideID] = append(h.clients[c.RideID], c)
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	clients := h.clients[c.RideID]
	for i, client := range clients {
		if client == c {
			h.clients[c.RideID] = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

func (h *Hub) Broadcast(rideID string, msg interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	payload, _ := json.Marshal(msg)
	clients := h.clients[rideID]
	for _, client := range clients {
		client.Conn.WriteMessage(websocket.TextMessage, payload)
	}
}
