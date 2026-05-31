package ws

type broadcastMsg struct {
	roomID  uint
	payload []byte
}

type onlineUsersRequest struct {
	roomID uint
	resp   chan []uint
}

type Hub struct {
	rooms          map[uint]map[*Client]bool
	onlineUsers    map[uint]map[uint]int
	Register       chan *Client
	Unregister     chan *Client
	broadcast      chan broadcastMsg
	onlineUsersReq chan onlineUsersRequest
}

func NewHub() *Hub {
	return &Hub{
		rooms:          make(map[uint]map[*Client]bool),
		onlineUsers:    make(map[uint]map[uint]int),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		broadcast:      make(chan broadcastMsg, 256),
		onlineUsersReq: make(chan onlineUsersRequest),
	}
}

func (h *Hub) GetOnlineUsers(roomID uint) []uint {
	resp := make(chan []uint, 1)
	h.onlineUsersReq <- onlineUsersRequest{roomID: roomID, resp: resp}
	return <-resp
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if h.rooms[client.RoomID] == nil {
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.rooms[client.RoomID][client] = true

			if h.onlineUsers[client.RoomID] == nil {
				h.onlineUsers[client.RoomID] = make(map[uint]int)
			}
			h.onlineUsers[client.RoomID][client.UserID]++

			client.send <- marshalEvent(EventOnlineUsers, onlineUsersData{UserIDs: h.getOnlineUserIDs(client.RoomID)})

			if h.onlineUsers[client.RoomID][client.UserID] == 1 {
				h.broadcastToRoom(client.RoomID, marshalEvent(EventUserJoined, userStatusData{UserID: client.UserID}))
			}

		case client := <-h.Unregister:
			if room := h.rooms[client.RoomID]; room != nil {
				if _, ok := room[client]; ok {
					delete(room, client)
					close(client.send)

					h.onlineUsers[client.RoomID][client.UserID]--

					if h.onlineUsers[client.RoomID][client.UserID] == 0 {
						delete(h.onlineUsers[client.RoomID], client.UserID)
						h.broadcastToRoom(client.RoomID, marshalEvent(EventUserLeft, userStatusData{UserID: client.UserID}))
					}
				}
			}

		case msg := <-h.broadcast:
			h.broadcastToRoom(msg.roomID, msg.payload)

		case req := <-h.onlineUsersReq:
			req.resp <- h.getOnlineUserIDs(req.roomID)
		}
	}
}

func (h *Hub) broadcastToRoom(roomID uint, payload []byte) {
	for client := range h.rooms[roomID] {
		select {
		case client.send <- payload:
		default:
			delete(h.rooms[roomID], client)
			close(client.send)
		}
	}
}

func (h *Hub) getOnlineUserIDs(roomID uint) []uint {
	userIDs := make([]uint, 0, len(h.onlineUsers[roomID]))
	for uid := range h.onlineUsers[roomID] {
		userIDs = append(userIDs, uid)
	}
	return userIDs
}
