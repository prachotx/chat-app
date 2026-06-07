package ws

import "encoding/json"

type WSEvent struct {
	Type string `json:"type"`
	Data any    `json:"data,omitempty"`
}

const (
	EventMessage     = "message"
	EventUserJoined  = "user_joined"
	EventUserLeft    = "user_left"
	EventOnlineUsers = "online_users"
)

type userStatusData struct {
	UserID uint `json:"user_id"`
}

type onlineUsersData struct {
	UserIDs []uint `json:"user_ids"`
}

func marshalEvent(eventType string, data any) []byte {
	b, _ := json.Marshal(WSEvent{Type: eventType, Data: data})
	return b
}

func MarshalEvent(eventType string, data any) ([]byte, error) {
	return json.Marshal(WSEvent{Type: eventType, Data: data})
}
