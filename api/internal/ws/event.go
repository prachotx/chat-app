package hub

import (
	"encoding/json"

	"github.com/prachotx/real-time-chat/api/internal/dto"
)

type IncomingMessage struct {
	Content string `json:"content"`
}

type Event struct {
	Type string `json:"type"`
	Data any    `json:"data,omitempty"`
}

func NewMessageEvent(msg dto.MessageResponse) ([]byte, error) {
	return json.Marshal(Event{Type: "message", Data: msg})
}
