package handler

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/prachotx/chat-app/api/config"
	"github.com/prachotx/chat-app/api/internal/dto"
	"github.com/prachotx/chat-app/api/internal/service"
	"github.com/prachotx/chat-app/api/pkg/response"
	"gorm.io/gorm"
)

type AuthHandler struct {
	authService service.AuthService
	cfg         *config.Config
}

func NewAuthHandler(authService service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{authService, cfg}
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var input dto.LoginDto
	if err := c.Bind().Body(&input); err != nil {
		return response.Send(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	tokenString, err := h.authService.Login(input)
	if err != nil {
		return response.Send(c, fiber.StatusUnauthorized, "Invalid credentials", nil)
	}

	cookie := fiber.Cookie{
		Name:     "access_token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   h.cfg.AppEnv == "production",
		SameSite: func() string {
			if h.cfg.AppEnv == "production" {
				return fiber.CookieSameSiteNoneMode
			}
			return fiber.CookieSameSiteStrictMode
		}(),
	}

	c.Cookie(&cookie)

	return response.Send(c, fiber.StatusOK, "Login successful", nil)
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req dto.RegisterDto
	if err := c.Bind().Body(&req); err != nil {
		return response.Send(c, fiber.StatusUnprocessableEntity, err.Error(), nil)
	}

	err := h.authService.Register(req)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return response.Send(c, fiber.StatusConflict, "Email already exists", nil)
		}
		return response.Send(c, fiber.StatusInternalServerError, "Failed to register user", nil)
	}

	return response.Send(c, fiber.StatusOK, "Register successful", nil)
}

func (h *AuthHandler) Profile(c fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	profile, err := h.authService.FindProfile(userID)
	if err != nil {
		return response.Send(c, fiber.StatusNotFound, "Profile not found", nil)
	}

	return response.Send(c, fiber.StatusOK, "Profile found", profile)
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   h.cfg.AppEnv == "production",
		SameSite: func() string {
			if h.cfg.AppEnv == "production" {
				return fiber.CookieSameSiteNoneMode
			}
			return fiber.CookieSameSiteStrictMode
		}(),
	}

	c.Cookie(&cookie)

	return response.Send(c, fiber.StatusOK, "Logout successful", nil)
}
