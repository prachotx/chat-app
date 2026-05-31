package ws

import "encoding/json"

// WSEvent คือ envelope สำหรับทุก message ที่ส่งออก WebSocket
// frontend ใช้ field "type" เพื่อแยกว่าต้องทำอะไรกับ message นี้
type WSEvent struct {
	Type string `json:"type"`
	Data any    `json:"data,omitempty"`
}

// event types ที่ส่งออก WebSocket
const (
	EventMessage     = "message"      // chat message ใหม่
	EventUserJoined  = "user_joined"  // มีคนเข้าห้อง
	EventUserLeft    = "user_left"    // มีคนออกจากห้อง
	EventOnlineUsers = "online_users" // list ผู้ใช้ที่ online อยู่ตอนนี้ (ส่งครั้งแรกตอน connect)
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

// MarshalEvent เหมือน marshalEvent แต่ export สำหรับใช้ใน package handler
func MarshalEvent(eventType string, data any) ([]byte, error) {
	return json.Marshal(WSEvent{Type: eventType, Data: data})
}
