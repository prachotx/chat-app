package main

import (
	"fmt"

	"github.com/fasthttp/websocket"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/prachotx/real-time-chat/api/config"
	"github.com/prachotx/real-time-chat/api/internal/handler"
	"github.com/prachotx/real-time-chat/api/internal/middleware"
	"github.com/prachotx/real-time-chat/api/internal/model"
	"github.com/prachotx/real-time-chat/api/internal/repository"
	"github.com/prachotx/real-time-chat/api/internal/service"
	"github.com/valyala/fasthttp"
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
	authHandler := handler.NewAuthHandler(authService)

	roomRepo := repository.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo)
	roomHandler := handler.NewRoomHandler(roomService)

	messageRepo := repository.NewMessageRepository(db)
	messageService := service.NewMessageService(messageRepo)
	messageHandler := handler.NewMessageHandler(messageService)

	api := app.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.Post("/login", authHandler.Login)
			auth.Post("/register", authHandler.Register)
			auth.Get("/me", middleware.AuthMiddleware, authHandler.Profile)
		}
		room := api.Group("/rooms")
		{
			room.Post("/", middleware.AuthMiddleware, roomHandler.Create)
			room.Get("/", middleware.AuthMiddleware, roomHandler.FindAll)
			room.Post("/:id/messages", middleware.AuthMiddleware, messageHandler.Create)
			room.Get("/:id/messages", middleware.AuthMiddleware, messageHandler.FindByRoomID)
		}
	}

	upgrader := websocket.FastHTTPUpgrader{
		CheckOrigin: func(ctx *fasthttp.RequestCtx) bool { return true },
	}

	app.Get("/ws/:roomId", func(c fiber.Ctx) error {
		roomId := c.Params("roomId")
		return upgrader.Upgrade(c.RequestCtx(), func(conn *websocket.Conn) {
			fmt.Printf("Connected: %s (room=%s)\n", conn.RemoteAddr(), roomId)
			defer fmt.Printf("Disconnected: %s (room=%s)\n", conn.RemoteAddr(), roomId)

			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					break
				}
			}
		})
	})

	app.Listen(":" + cfg.Port)
}
