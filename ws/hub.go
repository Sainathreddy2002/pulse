package ws

import (
	"errors"
	"sync"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[int64][]*Client
}

func NewWSConnection() *Hub {
	return &Hub{
		clients: make(map[int64][]*Client),
	}
}

func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client.userID] = append(h.clients[client.userID], client)
}

func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	conns := h.clients[client.userID]
	for i, c := range conns {
		if c == client {
			h.clients[client.userID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	if len(h.clients[client.userID]) == 0 {
		delete(h.clients, client.userID)
	}
}

func (h *Hub) Send(userID int64, messageType int, message []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	conns := h.clients[userID]
	if len(conns) == 0 {
		return errors.New("client not found")
	}

	alive := conns[:0]
	var lastErr error
	sent := 0
	for _, client := range conns {
		if err := client.conn.WriteMessage(messageType, message); err != nil {
			client.conn.Close()
			lastErr = err
			continue
		}
		alive = append(alive, client)
		sent++
	}

	if len(alive) == 0 {
		delete(h.clients, userID)
		if lastErr != nil {
			return lastErr
		}
		return errors.New("client not found")
	}
	h.clients[userID] = alive
	if sent == 0 {
		return lastErr
	}
	return nil
}
