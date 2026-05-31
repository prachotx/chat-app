package handler

import (
	"strconv"

	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/prachotx/real-time-chat/api/internal/dto"
	"github.com/prachotx/real-time-chat/api/internal/service"
	internalws "github.com/prachotx/real-time-chat/api/internal/ws"
	"github.com/prachotx/real-time-chat/api/pkg/response"
)

var upgrader = websocket.FastHTTPUpgrader{}

type WsHandler struct {
	hub            *internalws.Hub
	messageService service.MessageService
}

func NewWsHandler(hub *internalws.Hub, messageService service.MessageService) *WsHandler {
	return &WsHandler{hub, messageService}
}

func (h *WsHandler) Handle(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return fiber.ErrUnauthorized
	}

	roomID64, err := strconv.ParseUint(c.Params("room_id"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}
	roomID := uint(roomID64)

	return upgrader.Upgrade(c.RequestCtx(), func(conn *websocket.Conn) {
		client := internalws.NewClient(h.hub, conn, userID, roomID)
		h.hub.Register <- client

		go client.WritePump()
		client.ReadPump(func(content string) ([]byte, error) {
			resp, err := h.messageService.Save(dto.CreateMessageDto{Content: content}, roomID, userID)
			if err != nil {
				return nil, err
			}
			return internalws.MarshalEvent(internalws.EventMessage, resp)
		})
	})
}

func (h *WsHandler) GetOnlineUsers(c fiber.Ctx) error {
	roomID64, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	userIDs := h.hub.GetOnlineUsers(uint(roomID64))
	return response.Send(c, fiber.StatusOK, "Online users", userIDs)
}
