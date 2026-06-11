package repository

import (
	"github.com/prachotx/chat-app/api/internal/model"
	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(message *model.Message) error
	FindByRoomID(roomID uint) ([]model.Message, error)
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db}
}

func (r *messageRepository) Create(message *model.Message) error {
	if err := r.db.Create(message).Error; err != nil {
		return err
	}
	return r.db.Preload("Sender").First(message, message.ID).Error
}

func (r *messageRepository) FindByRoomID(roomID uint) ([]model.Message, error) {
	var messages []model.Message
	err := r.db.Preload("Sender").Where("room_id = ?", roomID).Find(&messages).Error
	return messages, err
}
