package middleware

import (
	"github.com/gofiber/fiber/v3"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/prachotx/chat-app/api/pkg/jwt"
	"github.com/prachotx/chat-app/api/pkg/response"
)

func AuthMiddleware(c fiber.Ctx) error {
	tokenString := c.Cookies("access_token")

	if tokenString == "" {
		return response.Send(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	token, err := jwt.ValidateToken(tokenString)

	if err != nil || !token.Valid {
		return response.Send(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	claims := token.Claims.(gojwt.MapClaims)

	userID := uint(claims["user_id"].(float64))

	c.Locals("user_id", userID)

	return c.Next()
}
