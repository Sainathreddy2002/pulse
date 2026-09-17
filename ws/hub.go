package ws

import (
	"errors"
	"sync"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[int64]*Client
}

func NewWSConnection() *Hub {
	return &Hub{
		clients: make(map[int64]*Client),
	}
}

func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client.userID] = client
}

func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[client.userID] == client {
		delete(h.clients, client.userID)
	}
}

func (h *Hub) Send(userID int64, messageType int, message []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	client, ok := h.clients[userID]
	if !ok {
		return errors.New("Client not found")
	}
	err := client.conn.WriteMessage(messageType, message)
	if err != nil {
		client.conn.Close()
		delete(h.clients, client.userID)
		return err
	}
	return nil
}
