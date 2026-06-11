package service

import (
	"github.com/prachotx/chat-app/api/internal/dto"
	"github.com/prachotx/chat-app/api/internal/model"
	"github.com/prachotx/chat-app/api/internal/repository"
	"github.com/prachotx/chat-app/api/pkg/crypto"
	"github.com/prachotx/chat-app/api/pkg/jwt"
)

type AuthService interface {
	Login(input dto.LoginDto) (string, error)
	Register(input dto.RegisterDto) error
	FindProfile(userID uint) (dto.UserResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo}
}

func (s *authService) Login(input dto.LoginDto) (string, error) {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return "", err
	}

	err = crypto.CheckPasswordHash(input.Password, user.Password)
	if err != nil {
		return "", err
	}

	tokenString, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) Register(input dto.RegisterDto) error {
	hashedPassword, err := crypto.HashPassword(input.Password)
	if err != nil {
		return err
	}

	user := &model.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
	}

	return s.userRepo.Create(user)
}

func (s *authService) FindProfile(userID uint) (dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
