package hub

import (
	"sync"

	"github.com/fasthttp/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	UserID uint
	mu     sync.Mutex
}

// Write is safe to call from multiple goroutines.
func (cl *Client) Write(messageType int, data []byte) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	return cl.Conn.WriteMessage(messageType, data)
}

type Hub struct {
	mu    sync.RWMutex
	rooms map[string]map[*Client]bool
}

func New() *Hub {
	return &Hub{
		rooms: make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Register(roomID string, cl *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*Client]bool)
	}
	h.rooms[roomID][cl] = true
}

func (h *Hub) Unregister(roomID string, cl *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	clients, ok := h.rooms[roomID]
	if !ok {
		return
	}
	delete(clients, cl)
	if len(clients) == 0 {
		delete(h.rooms, roomID)
	}
}

func (h *Hub) Count(roomID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.rooms[roomID])
}

func (h *Hub) Broadcast(roomID string, messageType int, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for cl := range h.rooms[roomID] {
		cl.Write(messageType, data)
	}
}
