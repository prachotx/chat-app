package service

import (
	"github.com/prachotx/real-time-chat/api/internal/dto"
	"github.com/prachotx/real-time-chat/api/internal/model"
	"github.com/prachotx/real-time-chat/api/internal/repository"
)

type MessageService interface {
	Create(input dto.CreateMessageDto, roomID uint, userID uint) error
	Save(input dto.CreateMessageDto, roomID uint, userID uint) (dto.MessageResponse, error)
	FindByRoomID(roomID uint) ([]dto.MessageResponse, error)
}

type messageService struct {
	messageRepo repository.MessageRepository
}

func NewMessageService(messageRepo repository.MessageRepository) MessageService {
	return &messageService{messageRepo}
}

func (s *messageService) Create(input dto.CreateMessageDto, roomID uint, userID uint) error {
	message := model.Message{
		Content:  input.Content,
		RoomID:   roomID,
		SenderID: userID,
	}

	return s.messageRepo.Create(&message)
}

func (s *messageService) Save(input dto.CreateMessageDto, roomID uint, userID uint) (dto.MessageResponse, error) {
	message := model.Message{
		Content:  input.Content,
		RoomID:   roomID,
		SenderID: userID,
	}

	if err := s.messageRepo.Create(&message); err != nil {
		return dto.MessageResponse{}, err
	}

	return dto.MessageResponse{
		ID:        message.ID,
		CreatedAt: message.CreatedAt.String(),
		Content:   message.Content,
		RoomID:    message.RoomID,
		SenderID:  message.SenderID,
		Sender: dto.UserResponse{
			ID:       message.Sender.ID,
			Username: message.Sender.Username,
			Email:    message.Sender.Email,
		},
	}, nil
}

func (s *messageService) FindByRoomID(roomID uint) ([]dto.MessageResponse, error) {
	messages, err := s.messageRepo.FindByRoomID(roomID)
	if err != nil {
		return nil, err
	}

	var result []dto.MessageResponse
	for _, m := range messages {
		result = append(result, dto.MessageResponse{
			ID:        m.ID,
			CreatedAt: m.CreatedAt.String(),
			UpdatedAt: m.UpdatedAt.String(),
			Content:   m.Content,
			RoomID:    m.RoomID,
			SenderID:  m.SenderID,
			Sender: dto.UserResponse{
				ID:       m.Sender.ID,
				Username: m.Sender.Username,
				Email:    m.Sender.Email,
			},
		})
	}
	return result, nil
}
