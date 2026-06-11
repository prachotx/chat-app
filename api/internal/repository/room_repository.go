package repository

import (
	"github.com/prachotx/chat-app/api/internal/model"
	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(room *model.Room) error
	FindAll() ([]model.Room, error)
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db}
}

func (r *roomRepository) Create(room *model.Room) error {
	return r.db.Create(room).Error
}

func (r *roomRepository) FindAll() ([]model.Room, error) {
	var rooms []model.Room
	err := r.db.Preload("CreatedBy").Find(&rooms).Error
	return rooms, err
}
