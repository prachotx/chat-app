package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/prachotx/chat-app/api/pkg/response"
)

var RateLimitMiddleware = limiter.New(limiter.Config{
	Max:        100,
	Expiration: 1 * time.Minute,
	LimitReached: func(c fiber.Ctx) error {
		return response.Send(c, fiber.StatusTooManyRequests, "Too many requests, please try again later", nil)
	},
})

var AuthRateLimitMiddleware = limiter.New(limiter.Config{
	Max:        10,
	Expiration: 1 * time.Minute,
	LimitReached: func(c fiber.Ctx) error {
		return response.Send(c, fiber.StatusTooManyRequests, "Too many requests, please try again later", nil)
	},
})
