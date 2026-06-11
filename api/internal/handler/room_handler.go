package handler

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/prachotx/chat-app/api/internal/dto"
	"github.com/prachotx/chat-app/api/internal/service"
	"github.com/prachotx/chat-app/api/pkg/response"
	"gorm.io/gorm"
)

type RoomHandler struct {
	roomService service.RoomService
}

func NewRoomHandler(roomService service.RoomService) *RoomHandler {
	return &RoomHandler{
		roomService,
	}
}

func (h *RoomHandler) Create(c fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req dto.CreateRoomDto
	if err := c.Bind().Body(&req); err != nil {
		return response.Send(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	err := h.roomService.Create(req, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return response.Send(c, fiber.StatusConflict, "Room name already exists", nil)
		}
		return response.Send(c, fiber.StatusInternalServerError, "Failed to create room", nil)
	}

	return response.Send(c, fiber.StatusCreated, "Room created successfully", nil)
}

func (h *RoomHandler) FindAll(c fiber.Ctx) error {
	rooms, err := h.roomService.FindAll()
	if err != nil {
		return response.Send(c, fiber.StatusInternalServerError, "Failed to retrieve rooms", nil)
	}

	return response.Send(c, fiber.StatusOK, "Rooms retrieved successfully", rooms)
}
