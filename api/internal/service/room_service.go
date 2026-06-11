package service

import (
	"github.com/prachotx/chat-app/api/internal/dto"
	"github.com/prachotx/chat-app/api/internal/model"
	"github.com/prachotx/chat-app/api/internal/repository"
)

type RoomService interface {
	Create(input dto.CreateRoomDto, userID uint) error
	FindAll() ([]dto.RoomResponse, error)
}

type roomService struct {
	roomRepo repository.RoomRepository
}

func NewRoomService(roomRepo repository.RoomRepository) RoomService {
	return &roomService{roomRepo}
}

func (s *roomService) Create(input dto.CreateRoomDto, userID uint) error {
	room := &model.Room{
		Name:        input.Name,
		CreatedByID: userID,
	}

	return s.roomRepo.Create(room)
}

func (s *roomService) FindAll() ([]dto.RoomResponse, error) {
	rooms, err := s.roomRepo.FindAll()
	if err != nil {
		return nil, err
	}

	result := make([]dto.RoomResponse, len(rooms))
	for i, room := range rooms {
		result[i] = dto.RoomResponse{
			ID:          room.ID,
			Name:        room.Name,
			CreatedByID: room.CreatedByID,
			CreatedBy: dto.RoomUserResponse{
				ID:       room.CreatedBy.ID,
				Username: room.CreatedBy.Username,
			},
		}
	}

	return result, nil
}
