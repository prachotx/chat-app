package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/prachotx/chat-app/api/config"
	"github.com/prachotx/chat-app/api/internal/handler"
	"github.com/prachotx/chat-app/api/internal/middleware"
	"github.com/prachotx/chat-app/api/internal/model"
	"github.com/prachotx/chat-app/api/internal/repository"
	"github.com/prachotx/chat-app/api/internal/service"
	internalws "github.com/prachotx/chat-app/api/internal/ws"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func main() {
	cfg := config.Load()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.User{}, &model.Room{}, &model.Message{})

	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{validate: validator.New()},
	})

	userRepo := repository.NewUserRepository(db)

	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService, cfg)

	roomRepo := repository.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo)
	roomHandler := handler.NewRoomHandler(roomService)

	messageRepo := repository.NewMessageRepository(db)
	messageService := service.NewMessageService(messageRepo)
	messageHandler := handler.NewMessageHandler(messageService)

	hub := internalws.NewHub()
	go hub.Run()
	wsHandler := handler.NewWsHandler(hub, messageService)

	api := app.Group("/api", middleware.RateLimitMiddleware)
	{
		auth := api.Group("/auth")
		{
			auth.Post("/login", middleware.AuthRateLimitMiddleware, authHandler.Login)
			auth.Post("/register", middleware.AuthRateLimitMiddleware, authHandler.Register)
			auth.Get("/me", middleware.AuthMiddleware, authHandler.Profile)
			auth.Get("/logout", middleware.AuthMiddleware, authHandler.Logout)
		}
		room := api.Group("/rooms")
		{
			room.Post("/", middleware.AuthMiddleware, roomHandler.Create)
			room.Get("/", middleware.AuthMiddleware, roomHandler.FindAll)
			room.Post("/:id/messages", middleware.AuthMiddleware, messageHandler.Create)
			room.Get("/:id/messages", middleware.AuthMiddleware, messageHandler.FindByRoomID)
		}
	}

	app.Get("/ws/:room_id", middleware.AuthMiddleware, wsHandler.Handle)

	app.Listen(":" + cfg.Port)
}
